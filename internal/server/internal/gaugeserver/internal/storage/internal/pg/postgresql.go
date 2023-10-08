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
