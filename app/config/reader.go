package config

import (
		"time"
	"github.com/WesJD/proxy-scraper/app/utils"
	"encoding/json"
	"io/ioutil"
)

type Configuration struct {
	Scraping ScrapingConfig `json:"scraping""`
	Influx InfluxConfig `json:"influx"`
	Checking CheckingConfig `json:"checking"`
}

type InfluxConfig struct {
	Address string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	UpdateEveryMs time.Duration `json:"updateEveryMs"`
}

type CheckingConfig struct {
	Services int `json:"services"`
	PerRound int `json:"perRound"`
}

type ScrapingConfig struct {
	DatabaseUrl string `json:"databaseUrl"`
	Static string `json:"static"`
	TimeoutMs int `json:"timeoutMs"`
}

var (
	Values *Configuration

	defaultConfig = &Configuration{
		Scraping: ScrapingConfig{
			DatabaseUrl: "",
			Static: "http://www.example.com",
			TimeoutMs: 1000,
		},
		Influx: InfluxConfig{
			Address: "http://localhost:8086",
			Username: "",
			Password: "",
			Database: "",
			UpdateEveryMs: 15000,
		},
		Checking: CheckingConfig{
			Services: 50,
			PerRound: 50,
		},
	}
)

func init() {
	path := utils.Resource("config.json")
	if !utils.Exists(path) {
		Values = defaultConfig
		encoded, err := json.MarshalIndent(config, "", "    ")
		utils.CheckError(err)
		err = ioutil.WriteFile(path, encoded, 0644)
		utils.CheckError(err)
	} else {
		Values = &Configuration{}
		encoded, err := ioutil.ReadFile(path)
		utils.CheckError(err)
		err = json.Unmarshal(encoded, config)
		utils.CheckError(err)
	}
	return
}