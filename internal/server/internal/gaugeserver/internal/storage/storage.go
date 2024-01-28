package storage

import (
	"context"
	"time"

	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/storage/internal/memory"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/storage/internal/pg"
)

// NewInMemoryStorage - создает хранилище в файле
func NewInMemoryStorage(
	ctx context.Context,
	filename string,
	restore bool,
	timeout time.Duration,
) (Storage, error) {
	return memory.NewStorage(ctx, filename, restore, timeout)
}

// NewPostgresqlStorage - создает хранилище postgresql
func NewPostgresqlStorage(dsn string) (Storage, error) {
	return pg.NewStorage(dsn)
}
