package main

import (
	"devops-tpl/internal/agent"
)

func main() {
	app := agent.NewHTTPClient()
	app.Run()
}
