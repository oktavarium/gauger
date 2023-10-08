package storage

import (
	"context"
	"io"

	"github.com/oktavarium/go-gauger/internal/shared"
)

type Storage interface {
	SaveGauge(context.Context, string, float64) error
	UpdateCounter(context.Context, string, int64) (int64, error)
	GetGauger(context.Context, string) (float64, bool)
	GetCounter(context.Context, string) (int64, bool)
	GetAll(context.Context) ([]byte, error)
	Ping(context.Context) error
	BatchUpdate(context.Context, io.Writer, []shared.Metric) error
}
