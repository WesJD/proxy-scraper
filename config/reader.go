package config

import (
	"encoding/json"
	"io/ioutil"
	"github.com/WesJD/proxy-scraper/utils"
)

func Read(path string, defaults interface{}) (value interface{}) {
	value = defaults
	if !utils.Exists(path) {
		encoded, err := json.MarshalIndent(value, "", "    ")
		utils.CheckError(err)
		err = ioutil.WriteFile(path, encoded, 0644)
		utils.CheckError(err)
	} else {
		encoded, err := ioutil.ReadFile(path)
		utils.CheckError(err)
		err = json.Unmarshal(encoded, value)
		utils.CheckError(err)
	}
	return
}
