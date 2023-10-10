package memory

import (
	"context"
	"fmt"
)

func (s *storage) SaveGauge(ctx context.Context, name string, val float64) error {
	s.gauge[name] = val
	if s.sync {
		return s.save()
	}
	return nil
}

func (s *storage) save() error {
	data, err := s.GetAll(context.Background())
	if err != nil {
		return fmt.Errorf("error on saving all: %w", err)
	}
	err = s.archive.Save(data)
	if err != nil {
		return fmt.Errorf("error on saving all: %w", err)
	}
	return nil
}
