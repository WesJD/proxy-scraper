package scraping

import (
	"time"
		"fmt"
		"github.com/WesJD/proxy-scraper/app/config"
	"github.com/WesJD/proxy-scraper/app/database"
	)

var (
	checkers = []Checker{
		&FreeProxyList{},
		&GetProxyList{},
		&Hidester{},
		&PremProxy{},
		&PubProxy{},
	}
)

type Checker interface {
	Check(url string, trueResponse string) (map[string]bool, error)
	WaitTime() time.Duration
}

func Start(config *config.Configuration, trueResponse string) {
	for _, checker := range checkers {
		go func(checker Checker) {
			for {
				proxies, err := checker.Check(config.Static, trueResponse)
				if err != nil {
					time.Sleep(checker.WaitTime())
					continue
				}
				fmt.Println(proxies)
				database.SubmitProxies(proxies)
				time.Sleep(checker.WaitTime())
			}
		}(checker)
	}

	return
}

