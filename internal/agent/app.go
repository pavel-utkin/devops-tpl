package agent

import (
	"devops-tpl/internal/agent/config"
	"devops-tpl/internal/agent/metricsuploader"
	"devops-tpl/internal/agent/statsreader"
	"log"
	"os"
	"syscall"
	"time"
)

type AppHTTP struct {
	isRun   bool
	timeLog struct {
		startTime       time.Time
		lastRefreshTime time.Time
		lastUploadTime  time.Time
	}
	metricsUplader *metricsuploader.MetricsUplader
	config         config.Config
}

func NewHTTPClient(config config.Config) *AppHTTP {
	var app AppHTTP
	app.config = config
	app.metricsUplader = metricsuploader.NewMetricsUploader(app.config.HTTPClientConnection, app.config.SignKey)

	return &app
}

func (app *AppHTTP) Run() {
	var metricsDump statsreader.MetricsDump
	signalChanel := make(chan os.Signal, 1)

	app.timeLog.startTime = time.Now()
	app.isRun = true

	tickerStatisticsRefresh := time.NewTicker(app.config.PollInterval)
	tickerStatisticsUpload := time.NewTicker(app.config.ReportInterval)

	for app.isRun {
		select {
		case timeTickerRefresh := <-tickerStatisticsRefresh.C:
			log.Println("Refresh")
			app.timeLog.lastRefreshTime = timeTickerRefresh
			metricsDump.Refresh()
		case timeTickerUpload := <-tickerStatisticsUpload.C:
			app.timeLog.lastUploadTime = timeTickerUpload
			log.Println("Upload")

			err := app.metricsUplader.MetricsUpload(metricsDump)
			if err != nil {
				log.Println("Error!")
				log.Println(err)

				app.Stop()
			}
		case osSignal := <-signalChanel:
			switch osSignal {
			case syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT:
				log.Println("syscall: " + osSignal.String())
			}
			app.Stop()
		}
	}
}

func (app *AppHTTP) Stop() {
	app.isRun = false
}

func (app *AppHTTP) IsRun() bool {
	return app.isRun
}
