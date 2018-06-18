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
)

func Initialize() {
	cfg := config.Read()

	httpclient.Defaults(httpclient.Map{
		httpclient.OPT_TIMEOUT:   7,
		httpclient.OPT_USERAGENT: "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:60.0) Gecko/20100101 Firefox/60.0",
	})

	database.Connect(cfg)
	scraping.Start(cfg)

	defer database.Influx.Close()
	defer database.Sql.Close()

	go func() {
		for {
			database.ReportStats()
			time.Sleep(1000 * 60 * 3 * time.Millisecond)
		}
	}()

	lock := true
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		<-signals
		lock = false
		os.Exit(0)
	}()

	for lock {
	}
}
