package agent

import (
	"errors"
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func loadConfig() (config, error) {
	var flagsConfig config
	flag.StringVar(&flagsConfig.Address, "a", "localhost:8080",
		"address and port of server's endpoint in notaion address:port")
	flag.IntVar(&flagsConfig.ReportInterval, "r", 10,
		"report interval in seconds")
	flag.IntVar(&flagsConfig.PollInterval, "p", 2,
		"poll interval in seconds")
	flag.Parse()

	if err := env.Parse(&flagsConfig); err != nil {
		return flagsConfig, fmt.Errorf("error on parsing env parameters: %w", err)
	}

	if len(flag.Args()) > 0 {
		return flagsConfig, errors.New("unrecognised flags")
	}

	flagsConfig.Address = "http://" + flagsConfig.Address

	return flagsConfig, nil
}
