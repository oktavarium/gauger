package main

import (
	"errors"
	"flag"

	"github.com/caarlos0/env/v6"
)

type config struct {
	Address         string `env:"ADDRESS"`
	Report_interval int    `env:"REPORT_INTERVAL"`
	Poll_interval   int    `env:"POLL_INTERVAL"`
}

var flagsConfig config

func parseFlags() error {
	flag.StringVar(&flagsConfig.Address, "a", "localhost:8080",
		"address and port of server's endpoint in notaion address:port")
	flag.IntVar(&flagsConfig.Report_interval, "r", 10,
		"report interval in seconds")
	flag.IntVar(&flagsConfig.Poll_interval, "p", 2,
		"apoll interval in seconds")
	flag.Parse()

	if err := env.Parse(&flagsConfig); err != nil {
		return err
	}

	if len(flag.Args()) > 0 {
		return errors.New("unrecognised flags")
	}

	flagsConfig.Address = "http://" + flagsConfig.Address

	return nil
}
