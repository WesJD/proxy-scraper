package main

import (
	"github.com/WesJD/proxy-scraper/config"
	"encoding/json"
)

type Configuration struct {
	Sql config.SQLDatabaseConfiguration `json:"sql"`
	Checking config.ProxyCheckerConfiguration `json:"checking"`
	HttpClient config.HttpClientDefaultsConfiguration `json:"client"`
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
