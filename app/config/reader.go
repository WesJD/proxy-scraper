package config

import (
	"github.com/WesJD/proxy-scraper/app/utils"
	"time"
	"encoding/json"
	"io/ioutil"
)

type Configuration struct {
	DatabaseUrl string `json:"databaseUrl"`
	Static string `json:"static"`
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
	OlderThan string `json:"olderThan"`
	EveryMs time.Duration `json:"everyMs"`
}

var defaultConfig = &Configuration{
	DatabaseUrl: "",
	Static: "http://www.example.com",
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
		OlderThan: "1 hour",
		EveryMs: 1000,
	},
}

func Read() (config *Configuration) {
	path := utils.Resource("config.json")
	if !utils.Exists(path) {
		config = defaultConfig
		encoded, err := json.Marshal(config)
		utils.CheckError(err)
		err = ioutil.WriteFile(path, encoded, 0644)
		utils.CheckError(err)
	} else {
		config = &Configuration{}
		encoded, err := ioutil.ReadFile(path)
		utils.CheckError(err)
		err = json.Unmarshal(encoded, config)
		utils.CheckError(err)
	}
	return
}