package server

import (
	"devops-tpl/internal/server/config"
	"devops-tpl/internal/server/middleware"
	"devops-tpl/internal/server/storage"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
)

type Server struct {
	storage   storage.MetricStorage
	chiRouter chi.Router
	config    config.Config
	startTime time.Time
}

func NewServer(config config.Config) *Server {
	log.Println(config)

	return &Server{
		config: config,
	}
}

func (server *Server) selectStorage() storage.MetricStorage {
	storageConfig := server.config.Store

	if storageConfig.DatabaseDSN != "" {
		log.Println("DB Storage")
		repository, err := storage.NewDBRepo(storageConfig)
		if err != nil {
			panic(err)
		}

		return repository
	}

	log.Println("Memory Storage")
	repository := storage.NewMetricsMemoryRepo(storageConfig)

	return repository
}

func (server *Server) initStorage() {
	metricsMemoryRepo := server.selectStorage()
	server.storage = metricsMemoryRepo

	if server.config.Store.Restore {
		server.storage.InitFromFile()
	}
}

func (server *Server) initRouter() {
	router := chi.NewRouter()

	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(middleware.GzipHandle)

	router.Get("/", server.PrintAllMetricStatic)
	router.Get("/ping", server.PingGetJSON)
	router.Get("/value/{statType}/{statName}", server.PrintMetricGet)

	router.Post("/value/", server.MetricValuePostJSON)
	router.Post("/updates/", server.UpdateMetricBatchJSON)

	router.Route("/update/", func(router chi.Router) {
		router.Post("/", server.UpdateMetricPostJSON)

		router.Post("/gauge/{statName}/{statValue}", server.UpdateGaugePost)
		router.Post("/counter/{statName}/{statValue}", server.UpdateCounterPost)
		router.Post("/{statType}/{statName}/{statValue}", server.UpdateNotImplementedPost)
	})

	server.chiRouter = router
}

func (server *Server) Run() {
	server.initStorage()
	defer server.storage.Close()
	server.initRouter()

	log.Fatal(http.ListenAndServe(server.config.ServerAddr, server.chiRouter))
}
