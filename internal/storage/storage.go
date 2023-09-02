package storage

type Storage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewStorage() *Storage {
	return &Storage{
		gauge:   map[string]float64{},
		counter: map[string]int64{},
	}
}

func (s *Storage) SaveGauge(name string, val float64) {
	s.gauge[name] = val
}

func (s *Storage) UpdateCounter(name string, val int64) {
	s.counter[name] += val
}

func (s *Storage) CheckMetricName(name string) bool {
	_, okGauge := s.gauge[name]
	_, okCounter := s.counter[name]
	return okCounter || okGauge
}
