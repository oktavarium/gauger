package main

import (
	"errors"
	"flag"

	"github.com/caarlos0/env/v6"
)

var flagRunAddr string

type config struct {
	flagRunAddr string `env:"ADDRESS"`
}

var flagsConfig config

func parseFlags() error {
	flag.StringVar(&flagsConfig.flagRunAddr, "a", "localhost:8080",
		"address and port of server in notaion address:port")
	flag.Parse()

	if err := env.Parse(&flagsConfig); err != nil {
		return err
	}

	if len(flag.Args()) > 0 {
		return errors.New("unrecognised flags")
	}

	return nil
}
