package pg

import (
	"context"
	"database/sql"
	"fmt"
)

func (s *storage) SaveGauge(ctx context.Context, name string, val float64) error {
	tx, err := s.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error occured on opening tx: %w", err)
	}
	defer tx.Rollback()

	var value float64
	row := tx.QueryRowContext(ctx, "SELECT value FROM gauge WHERE name = $1", name)
	err = row.Scan(&value)
	switch err {
	case sql.ErrNoRows:
		_, err = tx.ExecContext(ctx, "INSERT INTO gauge (name, value) VALUES ($1, $2)", name, val)
		if err != nil {
			return fmt.Errorf("error occured on inserting gauge: %w", err)
		}
	case nil:
		_, err = tx.ExecContext(ctx, "UPDATE gauge SET value = $2 WHERE name = $1", name, val)
		if err != nil {
			return fmt.Errorf("error occured on updating gauge: %w", err)
		}
	default:
		return fmt.Errorf("error occured on selecting gauge: %w", err)
	}

	return tx.Commit()
}
