package server

import (
	"devops-tpl/internal/server/config"
	"devops-tpl/internal/server/handlers"
	"devops-tpl/internal/server/storage"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"time"
)

type Server struct {
	startTime time.Time
	chiRouter chi.Router
}

func initRouter(memStatsStorage storage.MemStatsMemoryRepo) chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	//Маршруты
	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		handlers.PrintStatsValues(writer, request, memStatsStorage)
	})

	router.Route("/update", func(router chi.Router) {
		//Решил не делать 1 универсальный хэндлер, т.к. возможно в будущем будет разница
		//В обработке gauge и counter
		router.Post("/gauge/{statName}/{statValue}", func(writer http.ResponseWriter, request *http.Request) {
			handlers.UpdateGaugePost(writer, request, memStatsStorage)
		})
		router.Post("/counter/{statName}/{statValue}", func(writer http.ResponseWriter, request *http.Request) {
			handlers.UpdateCounterPost(writer, request, memStatsStorage)
		})
	})

	return router
}

func (server *Server) Run() {
	memStatsStorage := storage.NewMemStatsMemoryRepo()
	server.chiRouter = initRouter(memStatsStorage)

	fullHostAddr := fmt.Sprintf("%v:%v", config.Hostname, config.Port)
	log.Fatal(http.ListenAndServe(fullHostAddr, server.chiRouter))
}
