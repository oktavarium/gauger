package pg

import (
	"context"
	"database/sql"
	"fmt"
)

func (s *storage) UpdateCounter(ctx context.Context, name string, val int64) (int64, error) {
	tx, err := s.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("error occured on opening tx: %w", err)
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, "SELECT value FROM counter WHERE name = $1", name)
	var currentVal int64
	err = row.Scan(&currentVal)
	switch err {
	case sql.ErrNoRows:
		_, err = tx.ExecContext(ctx, "INSERT INTO counter (name, value) VALUES ($1, $2)", name, val)
		if err != nil {
			return 0, fmt.Errorf("error occured on inserting counter: %w", err)
		}
	case nil:
		_, err = tx.ExecContext(ctx, "UPDATE counter SET value = $2 WHERE name = $1", name, currentVal+val)
		if err != nil {
			return 0, fmt.Errorf("error occured on updating counter: %w", err)
		}
	default:
		return 0, fmt.Errorf("error occured on selecting counter: %w", err)
	}

	return val, tx.Commit()
}
