package main

import (
	"github.com/WesJD/proxy-scraper/utils"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"sync/atomic"
	"time"
	"github.com/WesJD/proxy-scraper/config"
	"database/sql"
	influxDB "github.com/influxdata/influxdb/client/v2"
)

const (
	updateWorkingProxy = "UPDATE proxies SET working = TRUE, checking = FALSE, consec_fails = 0 WHERE ip_port = ?;"
	updateNonWorkingProxy = "UPDATE proxies SET working = FALSE, checking = FALSE, consec_fails = consec_fails + 1 WHERE ip_port = ?;"
	callFormat = "CALL matchProxies(%d, NOW())"
)

var (
	defaultConfig = Configuration{
		Sql: config.DefaultSQLDatabaseConfiguration,
		Influx: config.DefaultInfluxDatabaseConfiguration,
		Reporting: config.DefaultStatisticsReportingConfiguration,
		HttpClient: config.DefaultHttpClientsDefaultConfiguration,
		Checking: config.DefaultProxyCheckerConfiguration,
		Instancing: CheckerConfiguration{
			Services: 50,
			PerRound: 50,
		},
	}

	preparedUpdateWorkingProxy *sql.Stmt
	preparedUpdateNonWorkingProxy *sql.Stmt

	influx influxDB.Client
	amountChecked int64
	lastReported int64

	cfg = readConfig()
)

func main() {
	cfg.HttpClient.Apply()

	var err error

	//influx
	influx, err = cfg.Influx.OpenConnection()
	utils.CheckError(err)

	// sql
	sql, err := cfg.Sql.OpenConnection()
	utils.CheckError(err)

	preparedUpdateWorkingProxy, err = sql.Prepare(updateWorkingProxy)
	utils.CheckError(err)

	preparedUpdateNonWorkingProxy, err = sql.Prepare(updateNonWorkingProxy)
	utils.CheckError(err)

	// the actual checking services
	for i := 0; i < cfg.Instancing.Services; i++ {
		go check(sql)
	}

	go reportStatistics()

	// lock until close
	select {
		case <-utils.WatchForKill():
			sql.Close()
			influx.Close()

			fmt.Println("Goodbye.")
			return
	}
}

func check(sql *sql.DB) {
	for {
		//cannot prepare a CALL statement... has to just stay here
		rows, err := sql.Query(fmt.Sprintf(callFormat, cfg.Instancing.PerRound))
		utils.CheckError(err)
		for rows.Next() {
			var ipPort string
			err = rows.Scan(&ipPort)
			utils.CheckError(err)

			checkResult := utils.CheckProxy(ipPort, cfg.Checking.StaticUrl)

			if checkResult {
				_, err = preparedUpdateWorkingProxy.Exec(ipPort)
			} else {
				_, err = preparedUpdateNonWorkingProxy.Exec(ipPort)
			}
			utils.CheckError(err)

			atomic.AddInt64(&amountChecked, 1)
		}
		rows.Close()
		time.Sleep(1)
	}
}

func reportStatistics() {
	batch, err := influxDB.NewBatchPoints(cfg.Reporting.GetBatchConfig(&cfg.Influx))
	utils.CheckError(err)

	currentTime := time.Now()
	secondsPassed := currentTime.Unix() - lastReported
	fmt.Println("sec", secondsPassed, "checked", amountChecked)

	if secondsPassed > 0 {
		fmt.Println("div", amountChecked / secondsPassed)

		fields := map[string]interface{}{
			"per second": amountChecked / secondsPassed,
		}

		point, err := influxDB.NewPoint("per_second", make(map[string]string), fields, currentTime)
		utils.CheckError(err)
		batch.AddPoint(point)

		lastReported = currentTime.Unix()
		amountChecked = 0
	}

	err = influx.Write(batch)
	utils.CheckError(err)

	time.Sleep(cfg.Reporting.Every)
}

func readConfig() (ret Configuration) {
	read := config.Read(utils.Resource("config.json"), defaultConfig)
	ret, ok := read.(Configuration)
	if !ok {
		log.Fatal("Couldn't cast config! watdehek")
	}
	return
}
