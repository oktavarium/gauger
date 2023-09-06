package server

import (
	"errors"
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type config struct {
	Address string `env:"ADDRESS"`
}

func loadConfig() (config, error) {
	var flagsConfig config
	flag.StringVar(&flagsConfig.Address, "a", "localhost:8080",
		"address and port of server in notaion address:port")
	flag.Parse()

	if err := env.Parse(&flagsConfig); err != nil {
		return flagsConfig, fmt.Errorf("error on parsing env parameters: %w", err)
	}

	if len(flag.Args()) > 0 {
		return flagsConfig, errors.New("unrecognised flags")
	}

	return flagsConfig, nil
}
