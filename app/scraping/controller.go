package scraping

import (
	"time"
	"github.com/ddliu/go-httpclient"
	"fmt"
	"github.com/WesJD/proxy-scraper/app/utils"
	"github.com/WesJD/proxy-scraper/app/config"
	"github.com/WesJD/proxy-scraper/app/database"
)

var (
	checkers = []Checker{
		GetProxyList{},
	}
)

type Checker interface {
	Check(url string, trueResponse string) (map[string]bool, error)
	WaitTime() time.Duration
}

func Start(config *config.Configuration) {
	res, err := httpclient.Get(config.Static)
	utils.CheckError(err)
	trueResponse, err := res.ToString()
	utils.CheckError(err)

	for _, checker := range checkers {
		go func() {
			for {
				proxies, err := checker.Check(config.Static, trueResponse)
				if err != nil {
					fmt.Println(":(", err.Error())
					continue
				}
				database.SubmitProxies(proxies)
				time.Sleep(checker.WaitTime())
			}
		}()
	}
}