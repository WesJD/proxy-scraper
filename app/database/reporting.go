package database

import (
	"github.com/WesJD/proxy-scraper/app/utils"
	"database/sql"
	"github.com/influxdata/influxdb/client/v2"
	"time"
	"fmt"
)

var (
	AmountChecked int64 = 0

	submitStatement *sql.Stmt
	getWorkingStatement *sql.Stmt
	lastReport = time.Now().Unix()
)

func makeStatements() {
	stmt, err := Sql.Prepare(Client, "insert-proxies")
	utils.CheckError(err)
	submitStatement = stmt

	stmt, err = Sql.Prepare(Client, "get-amount-working")
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
	rows.Close()

	batch, err := client.NewBatchPoints(batchConfig)
	utils.CheckError(err)

	fields := map[string]interface{}{
		"amount": amount,
	}
	point, err := client.NewPoint("proxy_count", make(map[string]string), fields, time.Now())
	utils.CheckError(err)
	batch.AddPoint(point)

	currentTime := time.Now()
	secondsPassed := currentTime.Unix() - lastReport
	fmt.Println("sec", secondsPassed, "checked", AmountChecked)
	if secondsPassed > 0 {
		fmt.Println("div", AmountChecked / secondsPassed)
		fields = map[string]interface{}{
			"per second": AmountChecked / secondsPassed,
		}
		point, err = client.NewPoint("per_second", make(map[string]string), fields, currentTime)
		utils.CheckError(err)
		batch.AddPoint(point)

		lastReport = currentTime.Unix()
		AmountChecked = 0
	}

	err = Influx.Write(batch)
	utils.CheckError(err)
}