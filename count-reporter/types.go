package main

import "github.com/WesJD/proxy-scraper/config"

type Configuration struct {
	Sql config.SQLDatabaseConfiguration `json:"sql"`
	Influx config.InfluxDatabaseConfiguration `json:"influx"`
	Reporting config.StatisticsReportingConfiguration `json:"reporting"`
}
