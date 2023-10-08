package memory

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/oktavarium/go-gauger/internal/shared"
)

func (s *storage) restore() error {
	data, err := s.archive.Restore()
	if err != nil {
		return fmt.Errorf("error on restoring archive: %w", err)
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		var metrics shared.Metric
		err := json.Unmarshal(scanner.Bytes(), &metrics)
		if err != nil {
			return fmt.Errorf("error on restoring archive: %w", err)
		}
		switch metrics.MType {
		case string(shared.GaugeType):
			s.SaveGauge(context.Background(), metrics.ID, *metrics.Value)
		case string(shared.CounterType):
			s.UpdateCounter(context.Background(), metrics.ID, *metrics.Delta)
		}
	}

	return nil
}
