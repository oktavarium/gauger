package storage

import (
	"time"

	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/storage/internal/memory"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/storage/internal/pg"
)

func NewInMemoryStorage(
	filename string,
	restore bool,
	timeout time.Duration,
) (Storage, error) {
	return memory.NewStorage(filename, restore, timeout)
}

func NewPostgresqlStorage(dsn string) (Storage, error) {
	return pg.NewStorage(dsn)
}
