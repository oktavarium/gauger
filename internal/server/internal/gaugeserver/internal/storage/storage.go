package storage

import (
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/storage/internal/memory"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/storage/internal/postgresql"
)

func NewInMemoryStorage(filename string, restore bool, timeout int) (Storage, error) {
	return memory.NewStorage(filename, restore, timeout)
}

func NewPostgresqlStorage(dsn string) (Storage, error) {
	return postgresql.NewStorage(dsn)
}
