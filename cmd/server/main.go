// HTTP сервер для runtime метрик
package main

import (
	"context"
	"devops-tpl/internal/server/config"
	"devops-tpl/internal/server/server"
	"errors"
	"fmt"
	"net/http"
)

func Profiling(addr string) {
	err := http.ListenAndServe(addr, nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
	}
}

func main() {
	config := config.LoadConfig()
	server := server.NewServer(config)

	if server.Config().ProfilingAddr != "" {
		go Profiling(server.Config().ProfilingAddr)
	}
	server.Run(context.Background())
}
