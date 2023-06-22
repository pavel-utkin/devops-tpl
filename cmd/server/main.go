// HTTP сервер для runtime метрик
package main

import (
	"context"
	"devops-tpl/internal/server/config"
	"devops-tpl/internal/server/server"
	"log"
	"net/http"
)

func Profiling(addr string) {
	log.Fatal(http.ListenAndServe(addr, nil))
}

func main() {
	config := config.LoadConfig()
	server := server.NewServer(config)

	if server.Config().ProfilingAddr != "" {
		go Profiling(server.Config().ProfilingAddr)
	}
	server.Run(context.Background())
}
