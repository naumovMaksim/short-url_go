package main

import (
	"net/http"

	"github.com/naumovMaksim/short-url_go/internal/handlers"
	"github.com/naumovMaksim/short-url_go/internal/service"
	"github.com/naumovMaksim/short-url_go/internal/storage"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	store := storage.NewMemoryStorage()
	serv := service.NewService(store)
	handler := handlers.NewHandler(serv)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.DefaultHandler)
	return http.ListenAndServe(`:8080`, mux)
}
