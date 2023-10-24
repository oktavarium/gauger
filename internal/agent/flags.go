package agent

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
)

type config struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	HashKey        string        `env:"KEY"`
	RateLimit      int           `env:"RATE_LIMIT"`
}

func loadConfig() (config, error) {
	var flagsConfig config
	flag.StringVar(&flagsConfig.Address, "a", "localhost:8080",
		"address and port of server's endpoint in notaion address:port")
	flag.DurationVar(&flagsConfig.ReportInterval, "r", 10,
		"report interval in seconds")
	flag.DurationVar(&flagsConfig.PollInterval, "p", 2,
		"poll interval in seconds")
	flag.StringVar(&flagsConfig.HashKey, "k", "",
		"key for hash")
	flag.IntVar(&flagsConfig.RateLimit, "l", 1,
		"requests limit")
	flag.Parse()

	if err := env.Parse(&flagsConfig); err != nil {
		return flagsConfig, fmt.Errorf("error on parsing env parameters: %w", err)
	}

	if len(flag.Args()) > 0 {
		return flagsConfig, errors.New("unrecognised flags")
	}

	if flagsConfig.RateLimit <= 0 {
		flagsConfig.RateLimit = 1
	}

	flagsConfig.Address = "http://" + flagsConfig.Address

	return flagsConfig, nil
}
