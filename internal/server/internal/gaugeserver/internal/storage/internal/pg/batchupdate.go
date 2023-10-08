package pg

import (
	"context"
	"fmt"

	"github.com/oktavarium/go-gauger/internal/shared"
)

func (s *storage) BatchUpdate(ctx context.Context, metrics []shared.Metric) error {
	tx, err := s.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error occured on creating tx on batchupdate: %w", err)
	}
	defer tx.Rollback()

	for _, v := range metrics {
		switch v.MType {
		case shared.GaugeType:
			if err := s.SaveGauge(ctx, v.ID, *v.Value); err != nil {
				return fmt.Errorf("error occured on saving gauge: %w", err)
			}
		case shared.CounterType:
			if _, err := s.UpdateCounter(ctx, v.ID, *v.Delta); err != nil {
				return fmt.Errorf("error occured on saving counter: %w", err)
			}
		}
	}
	return tx.Commit()
}
