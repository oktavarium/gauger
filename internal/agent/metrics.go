package agent

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"

	"github.com/go-resty/resty/v2"
	"github.com/oktavarium/go-gauger/internal/models"
)

const updatePath string = "update"

func NewMetrics() metrics {
	return metrics{
		gauges: gaugeMetrics{
			metrics: make(map[string]float64),
			mType:   models.GaugeType,
		},
		counters: counterMetrics{
			metrics: make(map[string]int64),
			mType:   models.CounterType,
		},
	}
}

func readMetrics(m *metrics) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	m.gauges.metrics["Alloc"] = float64(memStats.Alloc)
	m.gauges.metrics["TotalAlloc"] = float64(memStats.TotalAlloc)
	m.gauges.metrics["Sys"] = float64(memStats.Sys)
	m.gauges.metrics["Lookups"] = float64(memStats.Lookups)
	m.gauges.metrics["Frees"] = float64(memStats.Frees)
	m.gauges.metrics["Mallocs"] = float64(memStats.Mallocs)
	m.gauges.metrics["HeapAlloc"] = float64(memStats.HeapAlloc)
	m.gauges.metrics["HeapSys"] = float64(memStats.HeapSys)
	m.gauges.metrics["HeapIdle"] = float64(memStats.HeapIdle)
	m.gauges.metrics["HeapInuse"] = float64(memStats.HeapInuse)
	m.gauges.metrics["HeapReleased"] = float64(memStats.HeapReleased)
	m.gauges.metrics["HeapObjects"] = float64(memStats.HeapObjects)
	m.gauges.metrics["StackInuse"] = float64(memStats.StackInuse)
	m.gauges.metrics["StackSys"] = float64(memStats.StackSys)
	m.gauges.metrics["MSpanInuse"] = float64(memStats.MSpanInuse)
	m.gauges.metrics["MSpanSys"] = float64(memStats.MSpanSys)
	m.gauges.metrics["MCacheInuse"] = float64(memStats.MCacheInuse)
	m.gauges.metrics["MCacheSys"] = float64(memStats.MCacheSys)
	m.gauges.metrics["BuckHashSys"] = float64(memStats.BuckHashSys)
	m.gauges.metrics["GCSys"] = float64(memStats.GCSys)
	m.gauges.metrics["OtherSys"] = float64(memStats.OtherSys)
	m.gauges.metrics["NextGC"] = float64(memStats.NextGC)
	m.gauges.metrics["LastGC"] = float64(memStats.LastGC)
	m.gauges.metrics["PauseTotalNs"] = float64(memStats.PauseTotalNs)
	m.gauges.metrics["NumGC"] = float64(memStats.NumGC)
	m.gauges.metrics["NumForcedGC"] = float64(memStats.NumForcedGC)
	m.gauges.metrics["GCCPUFraction"] = float64(memStats.GCCPUFraction)
	m.gauges.metrics["RandomValue"] = rand.Float64()
	m.gauges.mType = models.GaugeType

	m.counters.metrics["PollCount"]++
	m.counters.mType = models.CounterType
}

func reportMetrics(address string, m *metrics) error {
	var metrics models.Metrics
	for k, v := range m.gauges.metrics {
		metrics.ID = k
		metrics.MType = string(models.GaugeType)
		metrics.Value = &v
		if err := makeUpdateRequest(fmt.Sprintf("%s/%s/", address, updatePath),
			metrics); err != nil {
			return fmt.Errorf("error on making update request: %w", err)
		}
	}

	metrics.Value = nil
	for k, v := range m.counters.metrics {
		metrics.ID = k
		metrics.MType = string(models.CounterType)
		metrics.Delta = &v
		if err := makeUpdateRequest(fmt.Sprintf("%s/%s/", address, updatePath),
			metrics); err != nil {
			return fmt.Errorf("error on making update request: %w", err)
		}
	}

	return nil
}

func makeUpdateRequest(endpoint string, metrics models.Metrics) error {
	var metricsResponse models.Metrics
	compressedMetrics, err := compressMetrics(metrics)
	if err != nil {
		return fmt.Errorf("error on compressing metrics on request: %w", err)
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(compressedMetrics).
		SetResult(&metricsResponse).
		Post(endpoint)

	if err != nil {
		return fmt.Errorf("error on making update request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.New("response status code is not OK (200)")
	}

	return nil
}
