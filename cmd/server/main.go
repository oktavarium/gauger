package main

import (
	"fmt"

	"github.com/oktavarium/go-gauger/internal/server"
	"github.com/oktavarium/go-gauger/internal/storage"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	err := parseFlags()
	if err != nil {
		return err
	}

	storage := storage.NewStorage()
	gs := server.NewGaugerServer(flagsConfig.Address, storage)

	fmt.Println("Running server on", flagsConfig.Address)

	return gs.ListenAndServe()
}
