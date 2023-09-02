package main

import (
	"log"

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
	storage := storage.NewStorage()
	gs := server.NewGaugerServer(srvAddr, storage)

	log.Println("Server started")

	return gs.ListenAndServe()
}
