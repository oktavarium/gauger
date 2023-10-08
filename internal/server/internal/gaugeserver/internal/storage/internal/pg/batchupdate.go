package pg

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/oktavarium/go-gauger/internal/shared"
)

func (s *storage) BatchUpdate(ctx context.Context, w io.Writer, metrics []shared.Metric) error {
	tx, err := s.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error occured on creating tx on batchupdate: %w", err)
	}
	defer tx.Rollback()
	encoder := json.NewEncoder(w)
	results := make([]shared.Metric, len(metrics))
	for _, v := range metrics {
		switch v.MType {
		case shared.GaugeType:
			if err := s.SaveGauge(ctx, v.ID, *v.Value); err != nil {
				return fmt.Errorf("error occured on saving gauge: %w", err)
			}
		case shared.CounterType:
			val, err := s.UpdateCounter(ctx, v.ID, *v.Delta)
			if err != nil {
				return fmt.Errorf("error occured on saving counter: %w", err)
			}
			*v.Delta = val
		}
		results = append(results, v)
	}
	if err := encoder.Encode(results[0]); err != nil {
		return fmt.Errorf("error occured on encoding result of batchupdate :%w", err)
	}
	return tx.Commit()
}
