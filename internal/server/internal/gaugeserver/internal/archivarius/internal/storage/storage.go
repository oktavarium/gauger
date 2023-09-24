package storage

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/oktavarium/go-gauger/internal/models"
)

type InMemoryStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemoryStorage() Storage {
	return Storage(&InMemoryStorage{
		gauge:   map[string]float64{},
		counter: map[string]int64{},
	})
}

func (s *InMemoryStorage) SaveGauge(name string, val float64) error {
	s.gauge[name] = val
	return nil
}

func (s *InMemoryStorage) UpdateCounter(name string, val int64) error {
	s.counter[name] += val
	return nil
}

func (s *InMemoryStorage) GetGauger(name string) (float64, bool) {
	val, ok := s.gauge[name]
	return val, ok
}

func (s *InMemoryStorage) GetCounter(name string) (int64, bool) {
	val, ok := s.counter[name]
	return val, ok
}

func (s *InMemoryStorage) GetAll() ([]byte, error) {
	allMetrics := make([]models.Metrics, 0, len(s.gauge)+len(s.counter))
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)

	for k, v := range s.gauge {
		var metrics models.Metrics

		metrics.ID = k
		metrics.MType = string(models.GaugeType)
		metrics.Value = &v
		allMetrics = append(allMetrics, metrics)

	}
	for k, v := range s.counter {
		var metrics models.Metrics

		metrics.ID = k
		metrics.MType = string(models.CounterType)
		metrics.Delta = &v

		allMetrics = append(allMetrics, metrics)
	}
	for _, v := range allMetrics {
		err := encoder.Encode(&v)
		if err != nil {
			return nil, fmt.Errorf("error on encoding data: %w", err)
		}

		// if i != len(allMetrics)-1 {
		// 	err = buffer.WriteByte('\n')
		// 	if err != nil {
		// 		return nil, fmt.Errorf("error on encoding data: %w", err)
		// 	}
		// }
	}
	return buffer.Bytes(), nil
}
