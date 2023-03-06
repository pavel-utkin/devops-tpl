package agent

import (
	"devops-tpl/internal/agent/metricsreader"
	"log"
	"os"
	"syscall"
	"time"

	"devops-tpl/internal/agent/config"
	"devops-tpl/internal/agent/requestorhandler"
	"github.com/go-resty/resty/v2"
)

type AppHTTP struct {
	isRun           bool
	startTime       time.Time
	lastRefreshTime time.Time
	lastUploadTime  time.Time
	client          *resty.Client
}

func (app *AppHTTP) initHTTPClient() {
	client := resty.New()

	client.
		SetRetryCount(config.ConfigClientRetryCount).
		SetRetryWaitTime(config.ConfigClientRetryWaitTime).
		SetRetryMaxWaitTime(config.ConfigClientRetryMaxWaitTime)

	app.client = client
}

func (app *AppHTTP) Run() {
	var memStatistics metricsreader.MemoryStatsDump
	signalChanel := make(chan os.Signal, 1)

	app.initHTTPClient()
	app.startTime = time.Now()
	app.isRun = true

	tickerStatisticsRefresh := time.NewTicker(config.ConfigPollInterval * time.Second)
	tickerStatisticsUpload := time.NewTicker(config.ConfigReportInterval * time.Second)

	for app.isRun {
		select {
		case timeTickerRefresh := <-tickerStatisticsRefresh.C:
			log.Println("Refresh")
			app.lastRefreshTime = timeTickerRefresh
			memStatistics.Refresh()
		case timeTickerUpload := <-tickerStatisticsUpload.C:
			app.lastUploadTime = timeTickerUpload
			log.Println("Upload")

			err := requestorhandler.MemoryStatsUpload(app.client, memStatistics)
			if err != nil {
				log.Println("Error!")
				log.Println(err)

				app.Stop()
			}
		case osSignal := <-signalChanel:
			switch osSignal {
			case syscall.SIGTERM:
				log.Println("syscall: SIGTERM")
			case syscall.SIGINT:
				log.Println("syscall: SIGINT")
			case syscall.SIGQUIT:
				log.Println("syscall: SIGQUIT")
			}
			app.Stop()
		}
	}
}

func (app *AppHTTP) Stop() {
	app.isRun = false
	os.Exit(1)
}

func (app *AppHTTP) IsRun() bool {
	return app.isRun
}
