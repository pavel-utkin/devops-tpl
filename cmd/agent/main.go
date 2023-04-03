package main

import (
	"devops-tpl/internal/agent"
	"devops-tpl/internal/agent/config"
)

func main() {
	config := config.LoadConfig()
	app := agent.NewHTTPClient(config)
	app.Run()
}
