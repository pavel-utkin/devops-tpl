// HTTP сервер для runtime метрик
package main

import (
	"context"
	"devops-tpl/internal/server/config"
	"devops-tpl/internal/server/server"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

func Profiling(addr string) {
	err := http.ListenAndServe(addr, nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
	}
}

func main() {

	ctx, ctxCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer ctxCancel()

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	config := config.LoadConfig()
	server := server.NewServer(config)

	if server.Config().ProfilingAddr != "" {
		go Profiling(server.Config().ProfilingAddr)
	}
	server.Run(ctx)
}
