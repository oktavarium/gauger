package archivarius

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/oktavarium/go-gauger/internal/models"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/archivarius/internal/archive"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/archivarius/internal/storage"
)

type archiver struct {
	storage.Storage
	archiveStorage archive.Archive
	timeout        int
}

func NewArchiver(filename string, restore bool, timeout int) (Archivarius, error) {
	ms := storage.NewMemoryStorage()
	as := archive.NewFileArchive(filename)
	ar := &archiver{ms, as, timeout}
	if restore {
		err := ar.restore()
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return nil, fmt.Errorf("error on restoring archive: %w", err)
			}

		}
	}
	if timeout != 0 {
		go func() {
			ticker := time.NewTicker(time.Duration(timeout) * time.Second)
			for range ticker.C {
				ar.save()

			}
		}()
	}
	return ar, nil
}

func (a *archiver) restore() error {
	data, err := a.archiveStorage.Restore()
	if err != nil {
		return fmt.Errorf("error on restoring archive: %w", err)
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		var metrics models.Metrics
		err := json.Unmarshal(scanner.Bytes(), &metrics)
		if err != nil {
			return fmt.Errorf("error on restoring archive: %w", err)
		}
		switch metrics.MType {
		case string(models.GaugeType):
			a.SaveGauge(metrics.ID, *metrics.Value)
		case string(models.CounterType):
			a.UpdateCounter(metrics.ID, *metrics.Delta)
		}
	}

	return nil
}

func (a *archiver) save() error {
	data, err := a.GetAll()
	if err != nil {
		return fmt.Errorf("error on saving all: %w", err)
	}
	err = a.archiveStorage.Save(data)
	if err != nil {
		return fmt.Errorf("error on saving all: %w", err)
	}
	return nil
}

func (a *archiver) SaveGauge(name string, val float64) error {
	a.Storage.SaveGauge(name, val)
	if a.timeout == 0 {
		return a.save()
	}
	return nil
}

func (a *archiver) UpdateCounter(name string, val int64) (int64, error) {
	val, err := a.Storage.UpdateCounter(name, val)
	if err != nil {
		return 0, fmt.Errorf("failed to update counter: %w", err)
	}
	if a.timeout == 0 {
		err := a.save()
		if err != nil {
			return 0, fmt.Errorf("failed to update counter: %w", err)
		}
	}
	return val, nil
}
