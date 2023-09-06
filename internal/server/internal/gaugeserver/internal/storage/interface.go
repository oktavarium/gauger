package storage

type Storage interface {
	SaveGauge(namse string, val float64)
	UpdateCounter(name string, val int64)
	GetGauger(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
	GetAll() string
}
