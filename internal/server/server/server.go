package server

import (
	"log"
	"net/http"
	"time"

	"devops-tpl/internal/server/config"
	"devops-tpl/internal/server/handlers"
	"devops-tpl/internal/server/middleware"
	"devops-tpl/internal/server/storage"
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

	log.Println("DB Storage LINE 33")
	log.Println("DB Storage LINE 33" + storageConfig.DatabaseDSN)
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
	if server.config.Store.Restore {
		repository.InitFromFile()
	}
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

	//Маршруты
	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		handlers.PrintStatsValues(writer, request, server.storage, server.config.TemplatesAbsPath)
	})

	router.Get("/ping", func(writer http.ResponseWriter, request *http.Request) {
		handlers.PingGet(writer, request, server.storage)
	})

	//json handler
	router.Post("/value/", func(writer http.ResponseWriter, request *http.Request) {
		handlers.JSONStatValue(writer, request, server.storage, server.config.SignKey)
	})

	router.Get("/value/{statType}/{statName}", func(writer http.ResponseWriter, request *http.Request) {
		handlers.PrintStatValue(writer, request, server.storage)
	})

	router.Route("/update/", func(router chi.Router) {
		//json handler
		router.Post("/", func(writer http.ResponseWriter, request *http.Request) {
			handlers.UpdateStatJSONPost(writer, request, server.storage, server.config.SignKey)
		})

		router.Post("/gauge/{statName}/{statValue}", func(writer http.ResponseWriter, request *http.Request) {
			handlers.UpdateGaugePost(writer, request, server.storage)
		})
		router.Post("/counter/{statName}/{statValue}", func(writer http.ResponseWriter, request *http.Request) {
			handlers.UpdateCounterPost(writer, request, server.storage)
		})
		router.Post("/{statType}/{statName}/{statValue}", func(writer http.ResponseWriter, request *http.Request) {
			handlers.UpdateNotImplementedPost(writer, request)
		})
	})

	server.chiRouter = router
}

func (server *Server) Run() {
	server.initStorage()
	defer server.storage.Close()
	server.initRouter()

	log.Fatal(http.ListenAndServe(server.config.ServerAddr, server.chiRouter))
}
