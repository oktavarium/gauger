package storage

type Storage interface {
	SaveGauge(string, float64) error
	UpdateCounter(string, int64) error
	GetGauger(string) (float64, bool)
	GetCounter(string) (int64, bool)
	GetAll() ([]byte, error)
}
