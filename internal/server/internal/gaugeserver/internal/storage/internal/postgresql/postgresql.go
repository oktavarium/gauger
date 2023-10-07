package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type storage struct {
	*sql.DB
}

func NewStorage(dsn string) (*storage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	s := &storage{
		db,
	}

	err = s.initDB(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error occured on init db when creating new storage: %w", err)
	}

	return s, nil
}

func (s *storage) initDB(ctx context.Context) error {
	err := s.Ping(ctx)
	if err != nil {
		return fmt.Errorf("error occured on db ping when initing db: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	_, err = s.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS gauge (id SERIAL PRIMARY KEY, name TEXT, value DOUBLE PRECISION)")
	if err != nil {
		return fmt.Errorf("error occured on creating table gauge: %w", err)
	}

	_, err = s.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS counter(id SERIAL PRIMARY KEY, name TEXT, value INTEGER)")
	if err != nil {
		return fmt.Errorf("error occured on creating table counter: %w", err)
	}

	return nil
}

func (s *storage) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	err := s.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("error occured on checking postgresql connection: %w", err)
	}

	return nil
}

func (s *storage) SaveGauge(name string, val float64) error {
	return nil
}

func (s *storage) UpdateCounter(ctx context.Context, name string, val int64) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	err := s.Ping(ctx)
	if err != nil {
		return 0, fmt.Errorf("error occured on db ping when updating counter: %w", err)
	}

	row := s.QueryRowContext(ctx, "SELECT name FROM counter WHERE name = $1", name)
	err = row.Scan()
	switch err {
	case sql.ErrNoRows:
		_, err = s.ExecContext(ctx, "INSERT INTO counter (name, value) VALUES ($1, $2)", name, val)
		if err != nil {
			return 0, fmt.Errorf("error occured in updating counter: %w", err)
		}
	case nil:
		_, err = s.ExecContext(ctx, "UPDATE counter SET value = $2 WHERE name = $1", name, val)
		if err != nil {
			return 0, fmt.Errorf("error occured in updating counter: %w", err)
		}
	default:
		return 0, fmt.Errorf("error occured in updating counter: %w", err)
	}

	return 0, nil
}

func (s *storage) GetGauger(name string) (float64, bool) {
	return 0.0, true
}

func (s *storage) GetCounter(name string) (int64, bool) {
	return 0, true
}

func (s *storage) GetAll() ([]byte, error) {

	return nil, nil
}
