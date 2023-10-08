package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/oktavarium/go-gauger/internal/shared"
)

func (s *storage) BatchUpdate(ctx context.Context, w io.Writer, metrics []shared.Metric) error {
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
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(metrics[0]); err != nil {
		return fmt.Errorf("error occured on encoding result of batchupdate :%w", err)
	}
	return nil
}
