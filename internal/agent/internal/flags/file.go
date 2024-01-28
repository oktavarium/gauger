package flags

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func (c *Config) parseConfigFile(file string) error {
	var config Config
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	if err := json.NewDecoder(r).Decode(&config); err != nil {
		return fmt.Errorf("error on decoding config: %w", err)
	}

	c.compare(config)

	return nil
}
