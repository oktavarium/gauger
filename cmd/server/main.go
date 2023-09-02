package main

import (
	"log"
	"net/http"

	"github.com/oktavarium/go-gauger/internal/handlers"
	"github.com/oktavarium/go-gauger/internal/server"
	"github.com/oktavarium/go-gauger/internal/storage"
)

const srvAddr string = "localhost:8080"

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	gs := server.NewGaugerServer(srvAddr)
	storage := storage.NewStorage()
	handlers := handlers.NewHandlers(storage)

	gs.Handle(`/`, http.HandlerFunc(handlers.UpdateHandle))

	log.Println("Server started")

	return gs.ListenAndServe()
}
