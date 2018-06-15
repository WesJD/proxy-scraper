package config

import (
	"runtime"
	"github.com/tkanos/gonfig"
	"path"
	"path/filepath"
	"github.com/WesJD/proxy-scraper/app/utils"
)

type Configuration struct {
	Static string
	Influx InfluxConfig
}

type InfluxConfig struct {
	Address string
	Username string
	Password string
	Database string
}

func Read() *Configuration {
	config := Configuration{}
	_, dirname, _, _ := runtime.Caller(0)
	err := gonfig.GetConf(path.Join(filepath.Dir(dirname), "config.json"), config)
	utils.CheckError(err)
	return &config
}