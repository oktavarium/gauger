package memory

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/storage/internal/memory/archive"
	"github.com/oktavarium/go-gauger/internal/shared"
)

type storage struct {
	gauge   map[string]float64
	counter map[string]int64
	archive archive.FileArchive
	sync    bool
}

func NewStorage(filename string, restore bool, timeout int) (*storage, error) {
	s := &storage{
		gauge:   map[string]float64{},
		counter: map[string]int64{},
		archive: archive.NewFileArchive(filename),
		sync:    timeout == 0,
	}

	if restore {
		err := s.restore()
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return nil, fmt.Errorf("failed to restore data from file: %w", err)
			}
		}
	}

	if !s.sync {
		go func() {
			ticker := time.NewTicker(time.Duration(timeout) * time.Second)
			for range ticker.C {
				s.save()
			}
		}()
	}

	return s, nil
}

func (s *storage) SaveGauge(name string, val float64) error {
	s.gauge[name] = val
	if s.sync {
		return s.save()
	}
	return nil
}

func (s *storage) UpdateCounter(name string, val int64) (int64, error) {
	s.counter[name] += val
	if s.sync {
		err := s.save()
		if err != nil {
			return 0, fmt.Errorf("failed to update counter: %w", err)
		}
	}
	return s.counter[name], nil
}

func (s *storage) GetGauger(name string) (float64, bool) {
	val, ok := s.gauge[name]
	return val, ok
}

func (s *storage) GetCounter(name string) (int64, bool) {
	val, ok := s.counter[name]
	return val, ok
}

func (s *storage) GetAll() ([]byte, error) {
	allMetrics := make([]shared.Metric, 0, len(s.gauge)+len(s.counter))
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)

	for k, v := range s.gauge {
		allMetrics = append(allMetrics, shared.NewGaugeMetric(k, &v))

	}
	for k, v := range s.counter {
		allMetrics = append(allMetrics, shared.NewCounterMetric(k, &v))
	}
	for _, v := range allMetrics {
		err := encoder.Encode(&v)
		if err != nil {
			return nil, fmt.Errorf("error on encoding data: %w", err)
		}
	}
	return buffer.Bytes(), nil
}

func (s *storage) restore() error {
	data, err := s.archive.Restore()
	if err != nil {
		return fmt.Errorf("error on restoring archive: %w", err)
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		var metrics shared.Metric
		err := json.Unmarshal(scanner.Bytes(), &metrics)
		if err != nil {
			return fmt.Errorf("error on restoring archive: %w", err)
		}
		switch metrics.MType {
		case string(shared.GaugeType):
			s.SaveGauge(metrics.ID, *metrics.Value)
		case string(shared.CounterType):
			s.UpdateCounter(metrics.ID, *metrics.Delta)
		}
	}

	return nil
}

func (s *storage) save() error {
	data, err := s.GetAll()
	if err != nil {
		return fmt.Errorf("error on saving all: %w", err)
	}
	err = s.archive.Save(data)
	if err != nil {
		return fmt.Errorf("error on saving all: %w", err)
	}
	return nil
}

func (s *storage) Ping() error {
	return nil
}
