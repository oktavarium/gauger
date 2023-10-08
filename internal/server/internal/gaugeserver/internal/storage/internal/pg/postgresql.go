package pg

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type storage struct {
	*sql.DB
}

func NewStorage(dsn string) (*storage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error occured on db ping when creating storage: %w", err)
	}

	s := &storage{
		db,
	}

	err = s.bootstrap(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error occured on init db when creating new storage: %w", err)
	}

	return s, nil
}

func (s *storage) Ping(ctx context.Context) error {
	err := s.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("error occured on checking postgresql connection: %w", err)
	}

	return nil
}

func (s *storage) SaveGauge(ctx context.Context, name string, val float64) error {
	tx, err := s.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error occured on opening tx: %w", err)
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, "SELECT value FROM gauge WHERE name = $1", name)
	err = row.Scan()
	switch err {
	case sql.ErrNoRows:
		_, err = tx.ExecContext(ctx, "INSERT INTO gauge (name, value) VALUES ($1, $2)", name, val)
		if err != nil {
			return fmt.Errorf("error occured on inserting gauge: %w", err)
		}
	case nil:
		_, err = tx.ExecContext(ctx, "UPDATE counter SET value = $2 WHERE name = $1", name, val)
		if err != nil {
			return fmt.Errorf("error occured on updating gauge: %w", err)
		}
	default:
		return fmt.Errorf("error occured on selecting gauge: %w", err)
	}

	return tx.Commit()
}

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

func (s *storage) GetGauger(ctx context.Context, name string) (float64, bool) {
	row := s.QueryRowContext(ctx, "SELECT value FROM gauge WHERE name = $1", name)
	var currentVal float64

	err := row.Scan(&currentVal)
	if err != nil {
		return 0.0, false
	}
	return currentVal, true
}

func (s *storage) GetCounter(ctx context.Context, name string) (int64, bool) {
	row := s.QueryRowContext(ctx, "SELECT value FROM counter WHERE name = $1", name)
	var currentVal int64

	err := row.Scan(&currentVal)
	if err != nil {
		return 0, false
	}
	return currentVal, true
}
