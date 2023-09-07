package storage

import (
	"fmt"
	"strconv"
	"strings"
)

type InMemoryStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewStorage() Storage {
	return Storage(&InMemoryStorage{
		gauge:   map[string]float64{},
		counter: map[string]int64{},
	})
}

func (s *InMemoryStorage) SaveGauge(name string, val float64) {
	s.gauge[name] = val
}

func (s *InMemoryStorage) UpdateCounter(name string, val int64) {
	s.counter[name] += val
}

func (s *InMemoryStorage) GetGauger(name string) (float64, bool) {
	val, ok := s.gauge[name]
	return val, ok
}

func (s *InMemoryStorage) GetCounter(name string) (int64, bool) {
	val, ok := s.counter[name]
	return val, ok
}

func (s *InMemoryStorage) GetAll() string {
	var sb strings.Builder
	for k, v := range s.gauge {
		sb.WriteString(fmt.Sprintf("%s: %s\n", k, strconv.FormatFloat(v, 'f', -1, 64)))
	}
	for k, v := range s.counter {
		sb.WriteString(fmt.Sprintf("%s: %d\n", k, v))
	}
	return sb.String()
}
