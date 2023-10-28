package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type storage struct {
	*pgxpool.Pool
}

func NewStorage(dsn string) (*storage, error) {
	ctx := context.Background()
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
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
