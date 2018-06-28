package database

import (
	"database/sql"
	"github.com/WesJD/proxy-scraper/app/utils"
	"github.com/WesJD/proxy-scraper/app/config"
	"github.com/influxdata/influxdb/client/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gchaincl/dotsql"
)

var (
	Client *sql.DB
	Influx client.Client
	Sql    *dotsql.DotSql
)

func Connect() {
	influx, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.Values.Influx.Address,
		Username: config.Values.Influx.Username,
		Password: config.Values.Influx.Password,
	})
	Influx = influx
	utils.CheckError(err)

	sqlDb, err := sql.Open("mysql", config.Values.Scraping.DatabaseUrl)
	Client = sqlDb
	utils.CheckError(err)

	setupDefaults()

	app, err := dotsql.LoadFromString(AppSql)
	utils.CheckError(err)
	Sql = app

	makeStatements()
}

func setupDefaults() {
	setup, err := dotsql.LoadFromString(SetupSql)
	utils.CheckError(err)

	//defaults
	names := []string {
		"setup-proxies",
	}
	for _, name := range names {
		_, err = setup.Exec(Client, name)
		utils.CheckError(err)
	}
}