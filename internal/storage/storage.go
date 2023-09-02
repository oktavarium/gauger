package storage

type MetricsSaver interface {
	SaveGauge(namse string, val float64)
	UpdateCounter(name string, val int64)
	CheckMetricName(name string) bool
}

type Storage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewStorage() *Storage {
	return &Storage{
		gauge: map[string]float64{
			"Alloc":         0.0,
			"TotalAlloc":    0.0,
			"Sys":           0.0,
			"Lookups":       0.0,
			"Frees":         0.0,
			"Mallocs":       0.0,
			"HeapAlloc":     0.0,
			"HeapSys":       0.0,
			"HeapIdle":      0.0,
			"HeapInuse":     0.0,
			"HeapReleased":  0.0,
			"HeapObjects":   0.0,
			"StackInuse":    0.0,
			"StackSys":      0.0,
			"MSpanInuse":    0.0,
			"MSpanSys":      0.0,
			"MCacheInuse":   0.0,
			"MCacheSys":     0.0,
			"BuckHashSys":   0.0,
			"GCSys":         0.0,
			"OtherSys":      0.0,
			"NextGC":        0.0,
			"LastGC":        0.0,
			"PauseTotalNs":  0.0,
			"NumGC":         0.0,
			"NumForcedGC":   0.0,
			"GCCPUFraction": 0.0,
			"RandomValue":   0.0,
		},
		counter: map[string]int64{
			"PollCount": 0,
		},
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
