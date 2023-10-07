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
	db, err := sql.Open("postgresql", dsn)
	if err != nil {
		return nil, err
	}

	s := &storage{
		db,
	}

	return s, nil
}

func (s *storage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
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
