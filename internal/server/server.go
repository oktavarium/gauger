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
	gs := gaugeserver.NewGaugerServer(flagsConfig.Address)

	return gs.ListenAndServe()
}
