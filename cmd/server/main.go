package main

import "devops-tpl/internal/server/server"

func main() {
	var httpServer server.Server
	httpServer.Run()
}
