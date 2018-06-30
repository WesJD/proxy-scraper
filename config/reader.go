package config

import (
	"encoding/json"
	"io/ioutil"
	"github.com/WesJD/proxy-scraper/utils"
		)

func Read(path string, defaults ReadableConfiguration) interface{} {
	if !utils.Exists(path) {
		encoded, err := json.MarshalIndent(&defaults, "", "    ")
		utils.CheckError(err)
		err = ioutil.WriteFile(path, encoded, 0644)
		utils.CheckError(err)
		return &defaults
	} else {
		encoded, err := ioutil.ReadFile(path)
		utils.CheckError(err)
		out, err := defaults.Read(encoded)
		utils.CheckError(err)
		return out
	}
}
