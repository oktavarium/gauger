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

	err = s.initDb(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error occured on init db when creating new storage: %w", err)
	}

	return s, nil
}

func (s *storage) initDb(ctx context.Context) error {
	err := s.Ping(ctx)
	if err != nil {
		return fmt.Errorf("error occured on db ping when initing db: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	_, err = s.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS gauge (name TEXT, value DOUBLE PRECISION)")
	if err != nil {
		return fmt.Errorf("error occured on creating table gauge: %w", err)
	}

	_, err = s.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS counter(name TEXT, value INTEGER)")
	if err != nil {
		return fmt.Errorf("error occured on creating table counter: %w", err)
	}

	// _, err = s.ExecContext(ctx, "GRANT ALL ON counter, gauge TO ")
	// if err != nil {
	// 	return fmt.Errorf("error occured on creating table counter: %w", err)
	// }
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

func (s *storage) UpdateCounter(name string, val int64) (int64, error) {
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
