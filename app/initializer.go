package app

import (
	"os"
	"os/signal"
	"syscall"
	"github.com/WesJD/proxy-scraper/app/config"
	"time"
	"github.com/WesJD/proxy-scraper/app/database"
	"github.com/WesJD/proxy-scraper/app/scraping"
	"github.com/ddliu/go-httpclient"
	"github.com/WesJD/proxy-scraper/app/checking"
	"github.com/WesJD/proxy-scraper/app/utils"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/WesJD/proxy-scraper/app/chrome"
	)

func Initialize() {
	cfg := config.Read()
	batchConfig := client.BatchPointsConfig{
		Database: cfg.Influx.Database,
		Precision: "s",
	}

	httpclient.Defaults(httpclient.Map{
		httpclient.OPT_TIMEOUT:   7,
		httpclient.OPT_USERAGENT: "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:60.0) Gecko/20100101 Firefox/60.0",
	})

	res, err := httpclient.Get(cfg.Static)
	utils.CheckError(err)
	trueResponse, err := res.ToString()
	utils.CheckError(err)

	database.Connect(cfg)
	scraping.Start(cfg, trueResponse)
	checking.Start(cfg, trueResponse)

	go func() {
		for {
			database.ReportStats(batchConfig)
			time.Sleep(cfg.Influx.UpdateEveryMs * time.Millisecond)
		}
	}()

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		<-signals

		_, err := database.Sql.Exec(database.Client, "set-all-not-checking")
		utils.CheckError(err)

		database.Influx.Close()
		database.Client.Close()
		chrome.CloseInstances()

		os.Exit(0)
	}()

	select {}
}
