package storage

import (
	"fmt"
	"strings"
)

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

func (s *Storage) GetGauger(name string) (float64, bool) {
	val, ok := s.gauge[name]
	return val, ok
}

func (s *Storage) GetCounter(name string) (int64, bool) {
	val, ok := s.counter[name]
	return val, ok
}

func (s *Storage) GetAll() string {
	var sb strings.Builder
	for k, v := range s.gauge {
		sb.WriteString(fmt.Sprintf("%q: %f\n", k, v))
	}
	for k, v := range s.counter {
		sb.WriteString(fmt.Sprintf("%q: %d\n", k, v))
	}
	return sb.String()
}
