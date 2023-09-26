package server

import (
	"fmt"

	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver"
	"github.com/oktavarium/go-gauger/internal/server/internal/logger"
)

func Run() error {
	flagsConfig, err := loadConfig()
	if err != nil {
		return fmt.Errorf("error on loading config: %w", err)
	}

	logger.Init(flagsConfig.LogLevel)

	gs, err := gaugeserver.NewGaugerServer(flagsConfig.Address,
		flagsConfig.FilePath,
		flagsConfig.Restore,
		flagsConfig.StoreInterval)
	if err != nil {
		return fmt.Errorf("error on creating gaugeserver: %w", err)
	}

	return gs.ListenAndServe()
}
