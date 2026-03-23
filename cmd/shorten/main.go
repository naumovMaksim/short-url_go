package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/naumovMaksim/short-url_go/internal/handlers"
	"github.com/naumovMaksim/short-url_go/internal/service"
	"github.com/naumovMaksim/short-url_go/internal/storage"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080", router()))
}

func router() http.Handler {
	store := storage.NewMemoryStorage()
	serv := service.NewService(store)
	handler := handlers.NewHandler(serv)

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", handler.GetHandler)
		r.Post("/", handler.AddHandler)
	})

	return r
}
