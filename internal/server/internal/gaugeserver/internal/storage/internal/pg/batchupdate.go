package pg

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/oktavarium/go-gauger/internal/shared"
)

var retry = 3
var delays = []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

func (s *storage) BatchUpdate(ctx context.Context, metrics []shared.Metric) error {
	for i := 0; ; i++ {
		err := s.batchUpdate(ctx, metrics)
		if err == nil || i >= retry {
			return err
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			select {
			case <-time.After(delays[i]):
				continue
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return err
	}
}

func (s *storage) batchUpdate(ctx context.Context, metrics []shared.Metric) error {
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
			val, err := s.UpdateCounter(ctx, v.ID, *v.Delta)
			if err != nil {
				return fmt.Errorf("error occured on saving counter: %w", err)
			}
			*v.Delta = val
		}
	}
	return tx.Commit()
}
