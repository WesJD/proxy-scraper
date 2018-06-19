package config

import (
	"github.com/tkanos/gonfig"
	"github.com/WesJD/proxy-scraper/app/utils"
	"time"
)

type Configuration struct {
	DatabaseUrl string
	Static string
	Influx InfluxConfig
	Checking CheckingConfig
}

type InfluxConfig struct {
	Address string
	Username string
	Password string
	Database string
	UpdateEveryMs time.Duration
}

type CheckingConfig struct {
	Services int
	PerRound int
	OlderThan string
	EveryMs time.Duration
}

func Read() *Configuration {
	config := &Configuration{}
	err := gonfig.GetConf(utils.Resource("config.json"), config)
	utils.CheckError(err)
	return config
}