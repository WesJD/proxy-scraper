package config

import (
	"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ddliu/go-httpclient"
	influxDB "github.com/influxdata/influxdb/client/v2"
)

type SQLDatabaseConfiguration struct {
	DatabaseUrl string `json:"databaseUrl"`
}

func (config *SQLDatabaseConfiguration) OpenConnection() (*sql.DB, error) {
	return sql.Open("sql", config.DatabaseUrl)
}

type InfluxDatabaseConfiguration struct {
	Address string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func (config *InfluxDatabaseConfiguration) OpenConnection() (influxDB.Client, error) {
	return influxDB.NewHTTPClient(influxDB.HTTPConfig{
		Addr:     config.Address,
		Username: config.Username,
		Password: config.Password,
	})
}

type HttpClientDefaultsConfiguration struct {
	UserAgent string `json:"userAgent"`
	Timeout int    `json:"timeoutMs"`
}

func (config *HttpClientDefaultsConfiguration) Apply() {
	httpclient.Defaults(httpclient.Map{
		httpclient.OPT_USERAGENT: config.UserAgent,
		httpclient.OPT_TIMEOUT_MS: config.Timeout,
	})
}

type ProxyCheckerConfiguration struct {
	StaticUrl string `json:"static"`
}

type StatisticsReportingConfiguration struct {
	Precision string `json:"precision"`
	Every time.Duration `json:"everyMs"`
}

func (config *StatisticsReportingConfiguration) GetBatchConfig(influx *InfluxDatabaseConfiguration) influxDB.BatchPointsConfig {
	return influxDB.BatchPointsConfig{
		Database: influx.Database,
		Precision: config.Precision,
	}
}
