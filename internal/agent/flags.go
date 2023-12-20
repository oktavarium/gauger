package agent

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
)

// config - структура хранения настроек агента
type config struct {
	Address        string        `env:"ADDRESS"`         // адрес сервиса сбора метрик
	ReportInterval time.Duration `env:"REPORT_INTERVAL"` // интервал отправки метрик
	PollInterval   time.Duration `env:"POLL_INTERVAL"`   // интервал сбора метрик
	HashKey        string        `env:"KEY"`             // ключ аутентификации
	RateLimit      int           `env:"RATE_LIMIT"`      // ограничение на количество поток
}

// loadConfig - загружает конфигурацию - из флагов и переменных окружения
func loadConfig() (config, error) {
	var flagsConfig config
	flag.StringVar(&flagsConfig.Address, "a", "localhost:8080",
		"address and port of server's endpoint in notaion address:port")
	flag.DurationVar(&flagsConfig.ReportInterval, "r", 10*time.Second,
		"report interval in seconds")
	flag.DurationVar(&flagsConfig.PollInterval, "p", 2*time.Second,
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
		fmt.Println(flag.Args())
		return flagsConfig, errors.New("unrecognised flags")
	}

	if flagsConfig.RateLimit <= 0 {
		flagsConfig.RateLimit = 1
	}

	flagsConfig.Address = "http://" + flagsConfig.Address

	return flagsConfig, nil
}
