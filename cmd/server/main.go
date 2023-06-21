package main

import (
	"devops-tpl/internal/server/config"
	"devops-tpl/internal/server/server"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func Profiling(addr string) {
	log.Fatal(http.ListenAndServe(addr, nil))
}

func main() {
	config := config.LoadConfig()
	server := server.NewServer(config)
	go Profiling(server.Config().ProfilingAddr)
	server.Run()
}
