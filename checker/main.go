package main

import (
	"github.com/WesJD/proxy-scraper/utils"
	"log"
	"fmt"
	"sync/atomic"
	"time"
	"github.com/WesJD/proxy-scraper/config"
	influxDB "github.com/influxdata/influxdb/client/v2"
	"database/sql"
	)

const (
	updateWorkingProxy = "UPDATE proxies SET working = TRUE, consec_fails = 0, consec_success = consec_success + 1 WHERE ip_port = ?;"
	updateNonWorkingProxy = "UPDATE proxies SET working = FALSE, consec_fails = consec_fails + 1, consec_success = 0 WHERE ip_port = ?;"
	callFormat = "CALL matchProxies(%d, NOW(), %d)"
	maxConsecFails = 50
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

	// client
	client, err := cfg.Sql.OpenConnection()
	utils.CheckError(err)

	preparedUpdateWorkingProxy, err = client.Prepare(updateWorkingProxy)
	utils.CheckError(err)

	preparedUpdateNonWorkingProxy, err = client.Prepare(updateNonWorkingProxy)
	utils.CheckError(err)

	// the actual checking services
	for i := 0; i < cfg.Instancing.Services; i++ {
		go check(client)
	}

	go reportStatistics()

	// lock until close
	select {
		case <-utils.WatchForKill():

			client.Close()
			influx.Close()

			fmt.Println("Goodbye.")
			return
	}
}

func check(sql *sql.DB) {
	for {
		//cannot prepare a CALL statement... has to just stay here
		rows, err := sql.Query(fmt.Sprintf(callFormat, cfg.Instancing.PerRound, maxConsecFails))
		utils.CheckError(err)
		for rows.Next() {
			var ipPort string
			err = rows.Scan(&ipPort)
			utils.CheckError(err)

			checkResult, checkError := utils.CheckProxyAndReason(ipPort, cfg.Checking.StaticUrl)

			if checkResult {
				_, err = preparedUpdateWorkingProxy.Exec(ipPort)
				fmt.Println("Got a working proxy: " + ipPort)
			} else {
				_, err = preparedUpdateNonWorkingProxy.Exec(ipPort)
				fmt.Println("Non working for reason " + checkError.Error())
			}
			utils.CheckError(err)

			atomic.AddInt64(&amountChecked, 1)
			time.Sleep(1 * time.Millisecond)
		}
		rows.Close()
		time.Sleep(50 * time.Millisecond)
	}
}

func reportStatistics() {
	for {
		batch, err := influxDB.NewBatchPoints(cfg.Reporting.GetBatchConfig(&cfg.Influx))
		utils.CheckError(err)

		currentTime := time.Now()
		secondsPassed := currentTime.Unix() - lastReported
		fmt.Println("sec", secondsPassed, "checked", amountChecked)

		if secondsPassed > 0 {
			fmt.Println("div", amountChecked/secondsPassed)

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

		time.Sleep(cfg.Reporting.Every * time.Millisecond) // uh, should this be using the config precision? because it's not right now
	}
}

func readConfig() (ret Configuration) {
	read := config.Read(utils.Resource("config.json"), defaultConfig)
	ret, ok := read.(Configuration)
	if !ok {
		log.Fatal("Couldn't cast config! watdehek")
	}
	return
}
