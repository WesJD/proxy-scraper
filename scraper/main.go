package main

import (
	"fmt"
	"reflect"
	"time"
	"github.com/WesJD/proxy-scraper/scraper/sites"
	"github.com/WesJD/proxy-scraper/utils"
	"os"
	"log"
	"github.com/WesJD/proxy-scraper/scraper/chrome"
	"github.com/WesJD/proxy-scraper/config"
	"database/sql"
)

const (
	createTables = `
		CREATE TABLE IF NOT EXISTS proxies (
  		ip_port      CHAR(40)         NOT NULL,
  		working      BOOL             NOT NULL,
  		last_checked TIMESTAMP        NOT NULL,
		consec_fails INTEGER UNSIGNED NOT NULL,
		consec_success INTEGER UNSIGNED NOT NULL,
  		UNIQUE (ip_port)
		);
	`
	createProxy = `
		INSERT INTO proxies (ip_port, working, consec_fails, consec_success) VALUES (?, ?, 0, 0)
		ON DUPLICATE KEY UPDATE working = VALUES(working), consec_fails = 0, consec_success = 0;
	`
)

var (
	checking = []sites.Site{
		&sites.FreeProxyList{},
		&sites.GetProxyList{},
		&sites.Hidester{},
		&sites.PremProxy{},
		&sites.ProxyNova{},
		&sites.PubProxy{},
	}
	defaultConfig = Configuration{
		Sql: config.DefaultSQLDatabaseConfiguration,
		Checking: config.DefaultProxyCheckerConfiguration,
		HttpClient: config.DefaultHttpClientsDefaultConfiguration,
	}

	preparedCreateProxy *sql.Stmt

	cfg = readConfig()
)

func main() {
	cfg.HttpClient.Apply()

	// sql
	sql, err := cfg.Sql.OpenConnection()
	utils.CheckError(err)

	_, err = sql.Exec(createTables)
	utils.CheckError(err)

	preparedCreateProxy, err = sql.Prepare(createProxy)
	utils.CheckError(err)

	// the actual scraping
	for _, site := range checking {
		go scrape(site)
	}

	// wait for kill
	select {
		case <-utils.WatchForKill():
			sql.Close()
			chrome.CloseInstances()

			fmt.Println("Goodbye.")
			os.Exit(0)
			return
	}
}

func scrape(site sites.Site) {
	for {
		proxies, err := site.Check(cfg.Checking.StaticUrl)
		if err != nil {
			fmt.Println(reflect.TypeOf(site), err)
			time.Sleep(site.WaitTime())
			continue
		}
		fmt.Println(reflect.TypeOf(site), proxies)

		for ipPort, working := range proxies {
			_, err = preparedCreateProxy.Exec(ipPort, working)
			utils.CheckError(err)
		}

		time.Sleep(site.WaitTime())
	}
}

func readConfig() (ret Configuration) {
	read := config.Read(utils.Resource("config.json"), defaultConfig)
	ret, ok := read.(Configuration)
	if !ok {
		log.Fatal("Couldn't cast config! watdehek")
	}
	return
}
