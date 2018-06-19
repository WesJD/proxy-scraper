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
	Sql             *sql.DB
	Influx          client.Client
	AppSql          *dotsql.DotSql
)

func Connect(config *config.Configuration) {
	influx, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.Influx.Address,
		Username: config.Influx.Username,
		Password: config.Influx.Password,
	})
	Influx = influx
	utils.CheckError(err)

	sqlDb, err := sql.Open("mysql", config.DatabaseUrl)
	Sql = sqlDb
	utils.CheckError(err)

	setupDefaults()

	app, err := dotsql.LoadFromFile(utils.Resource("app.sql"))
	utils.CheckError(err)
	AppSql = app

	makeStatements()
}

func setupDefaults() {
	setup, err := dotsql.LoadFromFile(utils.Resource("setup.sql"))
	utils.CheckError(err)

	//defaults
	names := []string {
		"setup-proxies",
	}
	for _, name := range names {
		_, err = setup.Exec(Sql, name)
		utils.CheckError(err)
	}
}