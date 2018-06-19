package database

import (
	"github.com/WesJD/proxy-scraper/app/utils"
	"database/sql"
	"github.com/influxdata/influxdb/client/v2"
		"time"
	"fmt"
)

var (
	submitStatement *sql.Stmt
	getWorkingStatement *sql.Stmt
)

func makeStatements() {
	stmt, err := AppSql.Prepare(Sql, "insert-proxies")
	utils.CheckError(err)
	submitStatement = stmt

	stmt, err = AppSql.Prepare(Sql, "get-amount-working")
	utils.CheckError(err)
	getWorkingStatement = stmt
}

func SubmitProxies(proxies map[string]bool) {
	for proxy, working := range proxies {
		_, err := submitStatement.Exec(proxy, working)
		utils.CheckError(err)
	}
}

func ReportStats(batchConfig client.BatchPointsConfig) {
	rows, err := getWorkingStatement.Query()
	utils.CheckError(err)

	var amount int
	rows.Next()
	err = rows.Scan(&amount)

	batch, err := client.NewBatchPoints(batchConfig)
	utils.CheckError(err)

	fields := map[string]interface{}{
		"amount": amount,
	}
	fmt.Println("amount", amount)
	point, err := client.NewPoint("proxy_count", make(map[string]string), fields, time.Now())
	utils.CheckError(err)
	batch.AddPoint(point)

	err = Influx.Write(batch)
	utils.CheckError(err)

	err = Influx.Close()
	utils.CheckError(err)
}