package main

import (
	"devops-tpl/internal/agent"
	"devops-tpl/internal/agent/config"
)

func main() {
	app := agent.NewHTTPClient(config.ClientRetryCount, config.ClientRetryWaitTime, config.ClientRetryMaxWaitTime)
	app.Run()
}
