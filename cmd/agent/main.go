package main

import (
	"devops-tpl/internal/agent"
	"devops-tpl/internal/agent/config"
)

func main() {
	app := agent.NewHTTPClient(config.AppConfig.HTTPClientConnection.RetryCount, config.AppConfig.HTTPClientConnection.RetryWaitTime, config.AppConfig.HTTPClientConnection.RetryMaxWaitTime)
	app.Run()
}
