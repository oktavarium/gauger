package main

import (
	"errors"
	"flag"

	"github.com/caarlos0/env/v6"
)

type config struct {
	flagEndpointAddr   string `env:"ADDRESS"`
	flagReportInterval int    `env:"REPORT_INTERVAL"`
	flagPollInterval   int    `env:"POLL_INTERVAL"`
}

var flagsConfig config

func parseFlags() error {
	flag.StringVar(&flagsConfig.flagEndpointAddr, "a", "localhost:8080",
		"address and port of server's endpoint in notaion address:port")
	flag.IntVar(&flagsConfig.flagReportInterval, "r", 10,
		"report interval in seconds")
	flag.IntVar(&flagsConfig.flagPollInterval, "p", 2,
		"apoll interval in seconds")
	flag.Parse()

	if err := env.Parse(&flagsConfig); err != nil {
		return err
	}

	if len(flag.Args()) > 0 {
		return errors.New("unrecognised flags")
	}

	flagsConfig.flagEndpointAddr = "http://" + flagsConfig.flagEndpointAddr

	return nil
}
