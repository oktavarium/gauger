package server

import (
	"errors"
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type config struct {
	Address       string `env:"ADDRESS"`
	LogLevel      string `env:"LOGLEVEL"`
	StoreInterval int    `env:"STORE_INTERVAL"`
	FilePath      string `env:"FILE_STORAGE_PATH"`
	Restore       bool   `env:"RESTORE"`
	DatabaseDSN   string `env:"DATABASE_DSN"`
}

func loadConfig() (config, error) {
	var flagsConfig config
	flag.StringVar(&flagsConfig.Address, "a", "localhost:8080",
		"address and port of server in notaion address:port")
	flag.StringVar(&flagsConfig.LogLevel, "l", "info",
		"log level")
	flag.IntVar(&flagsConfig.StoreInterval, "i", 300,
		"store interval")
	flag.StringVar(&flagsConfig.FilePath, "f", "/tmp/metrics-db.json",
		"file storage path")
	flag.BoolVar(&flagsConfig.Restore, "r", true,
		"restore metrics")
	flag.StringVar(&flagsConfig.DatabaseDSN, "d", "",
		"database connection string")
	flag.Parse()

	if err := env.Parse(&flagsConfig); err != nil {
		return flagsConfig, fmt.Errorf("error on parsing env parameters: %w", err)
	}

	if len(flag.Args()) > 0 {
		return flagsConfig, errors.New("unrecognised flags")
	}

	return flagsConfig, nil
}
