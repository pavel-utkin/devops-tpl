package agent

import (
	"devops-tpl/internal/agent/config"
	"devops-tpl/internal/agent/requesthandler"
	"devops-tpl/internal/agent/statsreader"
	"github.com/go-resty/resty/v2"
	"log"
	"os"
	"syscall"
	"time"
)

type AppHTTP struct {
	isRun           bool
	startTime       time.Time
	lastRefreshTime time.Time
	lastUploadTime  time.Time
	client          *resty.Client
}

func NewHTTPClient(clientRetryCount int, clientRetryWaitTime time.Duration, clientRetryMaxWaitTime time.Duration) *AppHTTP {
	var app AppHTTP
	client := resty.New()

	client.
		SetRetryCount(clientRetryCount).
		SetRetryWaitTime(clientRetryWaitTime).
		SetRetryMaxWaitTime(clientRetryMaxWaitTime)

	app.client = client
	return &app
}

func (app *AppHTTP) Run() {
	var memStatistics statsreader.MemoryStatsDump
	signalChanel := make(chan os.Signal, 1)

	app.startTime = time.Now()
	app.isRun = true

	tickerStatisticsRefresh := time.NewTicker(config.PollInterval * time.Second)
	tickerStatisticsUpload := time.NewTicker(config.ReportInterval * time.Second)

	for app.isRun {
		select {
		case timeTickerRefresh := <-tickerStatisticsRefresh.C:
			log.Println("Refresh")
			app.lastRefreshTime = timeTickerRefresh
			memStatistics.Refresh()
		case timeTickerUpload := <-tickerStatisticsUpload.C:
			app.lastUploadTime = timeTickerUpload
			log.Println("Upload")

			err := requesthandler.MemoryStatsUpload(app.client, memStatistics)
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
