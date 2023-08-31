package storage

type MetricsSaver interface {
	SaveGauge(namse string, val float64)
	UpdateCounter(name string, val int64)
}

type Storage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewStorage() *Storage {
	return &Storage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (s *Storage) SaveGauge(name string, val float64) {
	s.gauge[name] = val
}

func (s *Storage) UpdateCounter(name string, val int64) {
	s.counter[name] += val
}
