package server

import (
	"devops-tpl/internal/server/config"
	"devops-tpl/internal/server/handlers"
	"devops-tpl/internal/server/middleware"
	"devops-tpl/internal/server/storage"
	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"time"
)

type Server struct {
	chiRouter chi.Router
	config    config.Config
	startTime time.Time
}

func newRouter(metricsMemoryRepo storage.MetricStorage, templatesPath string) chi.Router {
	router := chi.NewRouter()

	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(middleware.GzipHandle)

	//Маршруты
	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		handlers.PrintStatsValues(writer, request, metricsMemoryRepo, templatesPath)
	})

	//json handler
	router.Post("/value/", func(writer http.ResponseWriter, request *http.Request) {
		handlers.JSONStatValue(writer, request, metricsMemoryRepo)
	})

	router.Get("/value/{statType}/{statName}", func(writer http.ResponseWriter, request *http.Request) {
		handlers.PrintStatValue(writer, request, metricsMemoryRepo)
	})

	router.Route("/update/", func(router chi.Router) {
		router.Post("/", func(writer http.ResponseWriter, request *http.Request) {
			handlers.UpdateStatJSONPost(writer, request, metricsMemoryRepo)
		})
		router.Post("/gauge/{statName}/{statValue}", func(writer http.ResponseWriter, request *http.Request) {
			handlers.UpdateGaugePost(writer, request, metricsMemoryRepo)
		})
		router.Post("/counter/{statName}/{statValue}", func(writer http.ResponseWriter, request *http.Request) {
			handlers.UpdateCounterPost(writer, request, metricsMemoryRepo)
		})
		router.Post("/{statType}/{statName}/{statValue}", func(writer http.ResponseWriter, request *http.Request) {
			handlers.UpdateNotImplementedPost(writer, request)
		})
	})

	return router
}

func NewServer(config config.Config) *Server {
	return &Server{
		config: config,
	}
}

func (server *Server) Run() {
	metricsMemoryRepo := storage.NewMetricsMemoryRepo(server.config.Store)
	defer metricsMemoryRepo.Close()
	if server.config.Store.Restore {
		metricsMemoryRepo.InitFromFile()
	}
	server.chiRouter = newRouter(metricsMemoryRepo, server.config.TemplatesAbsPath)

	log.Fatal(http.ListenAndServe(server.config.ServerAddr, server.chiRouter))
}
