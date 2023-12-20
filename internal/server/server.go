package server

import (
	"fmt"

	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver"
	"github.com/oktavarium/go-gauger/internal/server/internal/logger"
)

// Run - запускает сервис обработки метрик
func Run() error {
	fmt.Println("READ CONFIG!!!!!!!!!!!!!!!!!!!!")
	flagsConfig, err := loadConfig()
	if err != nil {
		return fmt.Errorf("error on loading config: %w", err)
	}
	fmt.Println("init logger!!!!!!!!!!!!!!!!!!!!")
	if err := logger.Init(flagsConfig.LogLevel); err != nil {
		return fmt.Errorf("error init logger: %w", err)
	}
	fmt.Println("NEW SERVER!!!!!!!!!!!!!!!!!!!!")
	gs, err := gaugeserver.NewGaugerServer(flagsConfig.Address,
		flagsConfig.FilePath,
		flagsConfig.Restore,
		flagsConfig.StoreInterval,
		flagsConfig.DatabaseDSN,
		flagsConfig.HashKey)
	if err != nil {
		return fmt.Errorf("error on creating gaugeserver: %w", err)
	}
	fmt.Println("START LISTEN!!!")
	return gs.ListenAndServe()
}
