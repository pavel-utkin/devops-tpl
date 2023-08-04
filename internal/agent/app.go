package agent

import (
	"context"
	"devops-tpl/internal/agent/config"
	"devops-tpl/internal/agent/metricsuploader"
	"devops-tpl/internal/agent/statsreader"
	"log"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type AppHTTP struct {
	isRun   bool
	timeLog struct {
		startTime       time.Time
		lastRefreshTime time.Time
		lastUploadTime  time.Time
	}
	metricsUplader      *metricsuploader.MetricsUplader
	metricsUploaderGRPC *metricsuploader.MetricsUploaderGRPC
	config              config.Config
}

func NewHTTPClient(config config.Config) *AppHTTP {
	var app AppHTTP
	app.config = config
	app.metricsUplader = metricsuploader.NewMetricsUploader(app.config.HTTPClientConnection, app.config.SignKey, app.config.PublicKeyRSA)

	if config.ServerGRPCAddr != "" {
		var err error
		app.metricsUploaderGRPC, err = metricsuploader.NewMetricsUploaderGRPC(app.config.ServerGRPCAddr)

		if err != nil {
			log.Fatal(err)
		}
	}

	return &app
}

func uploadMetrics(ctx context.Context, app *AppHTTP, metricsDump *statsreader.MetricsDump, wgRefresh *sync.WaitGroup) {
	wgRefresh.Wait()
	go func() {
		if app.metricsUploaderGRPC != nil {
			log.Println(app.metricsUploaderGRPC.Upload(ctx, *metricsDump))
			return
		}
		err := app.metricsUplader.MetricsUploadBatch(*metricsDump)
		if err != nil {
			log.Println(err)
		}
	}()
}

func (app *AppHTTP) Run(ctx context.Context) {
	metricsDump, err := statsreader.NewMetricsDump()
	if err != nil {
		log.Println(err)
		return
	}

	app.timeLog.startTime = time.Now()
	app.isRun = true

	tickerStatisticsRefresh := time.NewTicker(app.config.PollInterval)
	tickerStatisticsUpload := time.NewTicker(app.config.ReportInterval)
	wgRefresh := sync.WaitGroup{}

	workers := app.config.RateLimit

	errorGroup := new(errgroup.Group)
	for app.isRun {
		select {
		case timeTickerRefresh := <-tickerStatisticsRefresh.C:
			app.timeLog.lastRefreshTime = timeTickerRefresh
			wgRefresh.Add(2)

			for i := 0; i < workers; i++ {
				errorGroup.Go(func() error {
					go func() {
						defer wgRefresh.Done()
						metricsDump.Refresh()
					}()
					go func() {
						err = metricsDump.RefreshExtra()
						if err != nil {
							log.Println(err)
						}

						defer wgRefresh.Done()
					}()
					return nil
				})
			}
			errorGroup.Wait()
		case timeTickerUpload := <-tickerStatisticsUpload.C:
			app.timeLog.lastUploadTime = timeTickerUpload
			wgRefresh.Wait()

			for i := 0; i < workers; i++ {
				go func() {
					err = app.metricsUplader.MetricsUploadBatch(*metricsDump)
					if err != nil {
						log.Println("cant upload metrics ", err)
					}
				}()
			}
		case <-ctx.Done():
			uploadMetrics(ctx, app, metricsDump, &wgRefresh)
			wgRefresh.Wait()
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
