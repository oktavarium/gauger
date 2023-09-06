package server

import (
	"fmt"

	"github.com/oktavarium/go-gauger/internal/gaugeserver"
	"github.com/oktavarium/go-gauger/internal/handlers"
	"github.com/oktavarium/go-gauger/internal/storage"
)

func Run() error {
	flagsConfig, err := parseFlags()
	if err != nil {
		return fmt.Errorf("error on parsing flags: %w", err)
	}

	storage := storage.NewStorage()
	handler := handlers.NewHandler(storage)
	gs := gaugeserver.NewGaugerServer(flagsConfig.Address, handler)

	return gs.ListenAndServe()
}
