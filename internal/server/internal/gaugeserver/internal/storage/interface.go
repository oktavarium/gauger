package storage

import "context"

type Storage interface {
	SaveGauge(context.Context, string, float64) error
	UpdateCounter(context.Context, string, int64) (int64, error)
	GetGauger(context.Context, string) (float64, bool)
	GetCounter(context.Context, string) (int64, bool)
	GetAll(context.Context) ([]byte, error)
	Ping(context.Context) error
}
