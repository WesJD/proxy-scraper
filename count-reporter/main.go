package main

import (
	"github.com/WesJD/proxy-scraper/utils"
	"log"
	"github.com/WesJD/proxy-scraper/config"
	"time"
	"fmt"
	influxDB "github.com/influxdata/influxdb/client/v2"
)

const (
	getAmountWorking = "SELECT COUNT(*) FROM proxies WHERE working = TRUE;"
)

var (
	defaultConfig = Configuration{
		Sql: config.DefaultSQLDatabaseConfiguration,
		Influx: config.DefaultInfluxDatabaseConfiguration,
		Reporting: config.DefaultStatisticsReportingConfiguration,
	}

	cfg = readConfig()
)

func main() {
	// sql
	sql, err := cfg.Sql.OpenConnection()
	utils.CheckError(err)

	// influx
	influx, err := cfg.Influx.OpenConnection()
	utils.CheckError(err)

	// reporting
	for {
		rows, err := sql.Query(getAmountWorking)
		utils.CheckError(err)

		var amount int
		rows.Next()
		err = rows.Scan(&amount)
		rows.Close()

		batch, err := influxDB.NewBatchPoints(cfg.Reporting.GetBatchConfig(&cfg.Influx))
		utils.CheckError(err)

		fields := map[string]interface{}{
			"amount": amount,
		}

		point, err := influxDB.NewPoint("proxy_count", make(map[string]string), fields, time.Now())
		utils.CheckError(err)
		batch.AddPoint(point)

		err = influx.Write(batch)
		utils.CheckError(err)

		time.Sleep(cfg.Reporting.Every)
	}

	// wait for kill
	select {
		case <-utils.WatchForKill():
			sql.Close()
			influx.Close()

			fmt.Println("Goodbye.")
			return
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
