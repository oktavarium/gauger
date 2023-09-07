package server

import (
	"fmt"

	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver"
)

func Run() error {
	flagsConfig, err := loadConfig()
	if err != nil {
		return fmt.Errorf("error on loading config: %w", err)
	}

	gs := gaugeserver.NewGaugerServer(flagsConfig.Address)

	return gs.ListenAndServe()
}
