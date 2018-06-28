package main

import "github.com/WesJD/proxy-scraper/config"

type Configuration struct {
	Sql config.SQLDatabaseConfiguration `json:"sql"`
	Checking config.ProxyCheckerConfiguration `json:"checking"`
	HttpClient config.HttpClientDefaultsConfiguration `json:"client"`
}
