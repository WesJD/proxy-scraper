package database

import (
	"database/sql"
	"github.com/WesJD/proxy-scraper/app/utils"
	"github.com/WesJD/proxy-scraper/app/config"
	"github.com/influxdata/influxdb/client/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gchaincl/dotsql"
	"runtime"
	"path"
	"path/filepath"
)

var (
	Sql             *sql.DB
	Influx          client.Client
	submitStatement *sql.Stmt
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

	// Loads queries from file
	_, dirname, _, _ := runtime.Caller(0)
	dot, err := dotsql.LoadFromFile(path.Join(filepath.Dir(dirname), "setup.sql"))

	utils.CheckError(err)

	//defaults
	names := []string {
		"setup-proxies",
	}
	for _, name := range names {
		_, err = dot.Exec(Sql, name)
		utils.CheckError(err)
	}

	stmt, err := dot.Prepare(Sql, "insert-proxies")
	utils.CheckError(err)
	submitStatement = stmt
}
