package main

import "github.com/WesJD/proxy-scraper/config"

type Configuration struct {
	Sql config.SQLDatabaseConfiguration `json:"sql"`
	Influx config.InfluxDatabaseConfiguration `json:"influx"`
	Reporting config.StatisticsReportingConfiguration `json:"reporting"`
	HttpClient config.HttpClientDefaultsConfiguration `json:"client"`
	Checking config.ProxyCheckerConfiguration `json:"checking"`
	Instancing CheckerConfiguration `json:"instancing"`
}

type CheckerConfiguration struct {
	Services int `json:"services"`
	PerRound int `json:"perRound"`
}
