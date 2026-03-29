package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/naumovMaksim/short-url_go/internal/config"
	"github.com/naumovMaksim/short-url_go/internal/handlers"
	"github.com/naumovMaksim/short-url_go/internal/logger"
	"github.com/naumovMaksim/short-url_go/internal/service"
	"github.com/naumovMaksim/short-url_go/internal/storage"
)

func main() {
	conf := config.ParseFlags()
	handler := router(conf)
	logger.Initialize(conf.LogLevel)
	log.Fatal(http.ListenAndServe(conf.ServerAddress, handler))
}

func router(conf *config.Config) http.Handler {
	store := storage.NewMemoryStorage()
	serv := service.NewService(store, conf)
	handler := handlers.NewHandler(serv)

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Use(logger.RequestLogger)
		r.Get("/{id}", handler.GetHandler)
		r.Post("/", handler.AddHandler)
	})

	return r
}
