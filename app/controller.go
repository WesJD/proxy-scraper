package app

import (
	"os"
	"os/signal"
	"syscall"
	"github.com/ddliu/go-httpclient"
	"github.com/WesJD/proxy-scraper/app/utils"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/WesJD/proxy-scraper/app/config"
	"time"
	"sync"
	"fmt"
	"github.com/WesJD/proxy-scraper/app/checkers"
)

var (
	checkerStore = []checkers.Checker{
		checkers.PubProxy{},
	}
)

func Initialize() {
	cfg := config.Read()
	database, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     cfg.Influx.Address,
		Username: cfg.Influx.Username,
		Password: cfg.Influx.Password,
	})
	utils.CheckError(err)
	defer database.Close()

	httpclient.Defaults(httpclient.Map{
		httpclient.OPT_TIMEOUT: 7,
	})

	res, err := httpclient.Get(cfg.Static)
	utils.CheckError(err)
	trueResponse, err := res.ToString()
	utils.CheckError(err)
	fmt.Println("MY SITE", trueResponse)

	running := true

	wg := sync.WaitGroup{}

	batchConfig := client.BatchPointsConfig{
		Database:  cfg.Influx.Database,
		Precision: "s",
	}
	var batch client.BatchPoints
	wg.Add(1)
	go func() {
		defer wg.Done()

		for running {
			batch, err = client.NewBatchPoints(batchConfig)
			time.Sleep(1000 * 60 * 3 * time.Millisecond)
			utils.CheckError(database.Write(batch))
		}
	}()

	for _, value := range checkerStore {
		checker := value.(checkers.Checker)
		wg.Add(1)
		go func() {
			defer wg.Done()

			for running {
				result, err := checker.Check(cfg.Static, trueResponse)
				if err != nil {
					fmt.Println(":(", err.Error())
					continue
				}
				fmt.Println(result.WorkingProxies)

				fields := map[string]interface{}{
					"passing": result.Passing,
					"failing": result.Failing,
				}
				point, err := client.NewPoint("proxy", make(map[string]string), fields, time.Now())
				utils.CheckError(err)
				batch.AddPoint(point)

				time.Sleep(checker.WaitTime())
			}
		}()
	}

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		<-signals
		running = false

		utils.CheckError(database.Write(batch))
		utils.CheckError(database.Close())

		os.Exit(0)
	}()

	wg.Wait()
}