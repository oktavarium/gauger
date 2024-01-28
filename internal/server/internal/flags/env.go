package flags

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

func (c *Config) loadEnv() error {
	var config Config
	if err := env.Parse(&config); err != nil {
		return fmt.Errorf("error on parsing env parameters: %w", err)
	}

	c.compare(config)

	return nil
}
