package pg

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/oktavarium/go-gauger/internal/shared"
)

var retry = 3
var delays = []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

func (s *storage) BatchUpdate(ctx context.Context, metrics []shared.Metric) error {
	var gauge []shared.Metric
	var counter []shared.Metric
	for _, v := range metrics {
		switch v.MType {
		case shared.GaugeType:
			gauge = append(gauge, v)
		case shared.CounterType:
			counter = append(counter, v)
		}
	}

	for i := 0; ; i++ {
		err := s.batchUpdate(ctx, gauge, counter)
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

func (s *storage) batchUpdate(ctx context.Context, gauge []shared.Metric, counter []shared.Metric) error {
	tx, err := s.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("error occured on creating tx on batchupdate: %w", err)
	}
	defer tx.Rollback(ctx)

	gaugeBatch := pgx.Batch{}
	gaugeQuery := `
		INSERT INTO gauge (name, value) VALUES ($1, $2)
		ON CONFLICT (name) DO
		UPDATE SET value = $2
	`
	for _, v := range gauge {
		gaugeBatch.Queue(gaugeQuery, v.ID, v.Value)
	}

	err = s.SendBatch(ctx, &gaugeBatch).Close()
	if err != nil {
		return fmt.Errorf("error occured on making batch gauge update: %w", err)
	}

	counterBatch := pgx.Batch{}
	counterQuery := `
		INSERT INTO counter (name, value) VALUES ($1, $2)
		ON CONFLICT (name) DO
		UPDATE SET value = counter.value + $2
	`
	for _, v := range gauge {
		counterBatch.Queue(counterQuery, v.ID, v.Delta)
	}

	err = s.SendBatch(ctx, &counterBatch).Close()
	if err != nil {
		return fmt.Errorf("error occured on making batch counter update: %w", err)
	}

	return tx.Commit(ctx)
}
