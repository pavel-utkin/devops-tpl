package config

import "time"

const (
	ClientRetryCount       = 3
	ClientRetryWaitTime    = 10 * time.Second
	ClientRetryMaxWaitTime = 90 * time.Second
	PollInterval           = 2  //Seconds
	ReportInterval         = 10 //Seconds
	ServerHost             = "127.0.0.1"
	ServerPort             = 8080
)
