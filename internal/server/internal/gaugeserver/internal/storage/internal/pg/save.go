package pg

import (
	"context"
	"fmt"
)

func (s *storage) SaveGauge(ctx context.Context, name string, val float64) error {
	tx, err := s.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error occured on opening tx: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO gauge (name, value) VALUES ($1, $2)
		ON CONFLICT (name) DO
		UPDATE gauge SET value = gauge.value + %2`, name, val)

	if err != nil {
		return fmt.Errorf("error occured on inserting gauge: %w", err)
	}

	return tx.Commit()
}
