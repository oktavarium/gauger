package agent

import (
	"fmt"
	"time"
)

func Run() error {
	flagsConfig, err := loadConfig()
	if err != nil {
		return fmt.Errorf("error on loading config: %w", err)
	}

	metrics := NewMetrics()
	var sleepCounter int
	for {
		time.Sleep(1 * time.Second)
		sleepCounter++
		if sleepCounter%flagsConfig.PollInterval == 0 {
			readMetrics(&metrics)
		}
		if sleepCounter%flagsConfig.ReportInterval == 0 {
			if err := reportMetrics(flagsConfig.Address, &metrics); err != nil {
				//panic(fmt.Errorf("error reporting metrics: %w", err))
				continue
			}
		}
	}
}
