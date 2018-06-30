package main

import (
	"github.com/WesJD/proxy-scraper/config"
	"encoding/json"
	)

type Configuration struct {
	Sql config.SQLDatabaseConfiguration `json:"sql"`
	Influx config.InfluxDatabaseConfiguration `json:"influx"`
	Reporting config.StatisticsReportingConfiguration `json:"reporting"`
	HttpClient config.HttpClientDefaultsConfiguration `json:"client"`
	Checking config.ProxyCheckerConfiguration `json:"checking"`
	Instancing CheckerConfiguration `json:"instancing"`
}

func (config Configuration) Read(data []byte) (out interface{}, err error) {
	ret := Configuration{}
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return
	}
	out = ret
	return
}

type CheckerConfiguration struct {
	Services int `json:"services"`
	PerRound int `json:"perRound"`
}