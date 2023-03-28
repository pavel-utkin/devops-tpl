package main

import (
	"devops-tpl/internal/server/config"
	"devops-tpl/internal/server/server"
)

func main() {
	config := config.LoadConfig()
	server := server.NewServer(config)
	server.Run()
}
