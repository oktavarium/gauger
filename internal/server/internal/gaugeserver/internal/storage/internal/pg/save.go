package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *storage) SaveGauge(ctx context.Context, name string, val float64) error {
	tx, err := s.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("error occured on opening tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO gauge (name, value) VALUES ($1, $2)
		ON CONFLICT (name) DO
		UPDATE SET value = $2`, name, val)

	if err != nil {
		return fmt.Errorf("error occured on saving gauge: %w", err)
	}

	return tx.Commit(ctx)
}
