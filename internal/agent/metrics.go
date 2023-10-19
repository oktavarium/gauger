package agent

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"

	"github.com/go-resty/resty/v2"
	"github.com/oktavarium/go-gauger/internal/shared"
)

const updatePath string = "updates"

type metrics struct {
	gauges   map[string]float64
	counters map[string]int64
}

func NewMetrics() metrics {
	return metrics{
		make(map[string]float64),
		make(map[string]int64),
	}
}

func readMetrics(m *metrics) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	m.gauges["Alloc"] = float64(memStats.Alloc)
	m.gauges["TotalAlloc"] = float64(memStats.TotalAlloc)
	m.gauges["Sys"] = float64(memStats.Sys)
	m.gauges["Lookups"] = float64(memStats.Lookups)
	m.gauges["Frees"] = float64(memStats.Frees)
	m.gauges["Mallocs"] = float64(memStats.Mallocs)
	m.gauges["HeapAlloc"] = float64(memStats.HeapAlloc)
	m.gauges["HeapSys"] = float64(memStats.HeapSys)
	m.gauges["HeapIdle"] = float64(memStats.HeapIdle)
	m.gauges["HeapInuse"] = float64(memStats.HeapInuse)
	m.gauges["HeapReleased"] = float64(memStats.HeapReleased)
	m.gauges["HeapObjects"] = float64(memStats.HeapObjects)
	m.gauges["StackInuse"] = float64(memStats.StackInuse)
	m.gauges["StackSys"] = float64(memStats.StackSys)
	m.gauges["MSpanInuse"] = float64(memStats.MSpanInuse)
	m.gauges["MSpanSys"] = float64(memStats.MSpanSys)
	m.gauges["MCacheInuse"] = float64(memStats.MCacheInuse)
	m.gauges["MCacheSys"] = float64(memStats.MCacheSys)
	m.gauges["BuckHashSys"] = float64(memStats.BuckHashSys)
	m.gauges["GCSys"] = float64(memStats.GCSys)
	m.gauges["OtherSys"] = float64(memStats.OtherSys)
	m.gauges["NextGC"] = float64(memStats.NextGC)
	m.gauges["LastGC"] = float64(memStats.LastGC)
	m.gauges["PauseTotalNs"] = float64(memStats.PauseTotalNs)
	m.gauges["NumGC"] = float64(memStats.NumGC)
	m.gauges["NumForcedGC"] = float64(memStats.NumForcedGC)
	m.gauges["GCCPUFraction"] = float64(memStats.GCCPUFraction)
	m.gauges["RandomValue"] = rand.Float64()

	m.counters["PollCount"]++
}

func reportMetrics(address string, key string, m *metrics) error {
	allMetrics := make([]shared.Metric, 0, len(m.gauges)+len(m.counters))
	for k, v := range m.gauges {
		allMetrics = append(allMetrics, shared.NewGaugeMetric(k, &v))
	}

	for k, v := range m.counters {
		allMetrics = append(allMetrics, shared.NewCounterMetric(k, &v))
	}

	if err := makeBatchUpdateRequest(fmt.Sprintf("%s/%s/", address, updatePath), key, allMetrics); err != nil {
		return fmt.Errorf("error on making batchupdate request: %w", err)
	}

	return nil
}

func makeBatchUpdateRequest(endpoint string, key string, metrics []shared.Metric) error {
	var metricsResponse shared.Metric
	compressedMetrics, err := compressMetrics(metrics)
	if err != nil {
		return fmt.Errorf("error on compressing metrics on request: %w", err)
	}

	hashedMetrics, err := hashData([]byte(key), compressedMetrics)
	if err != nil {
		return fmt.Errorf("error occured in calculating hmac: %w", err)
	}
	fmt.Println(string(compressedMetrics), string(hashedMetrics))
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetHeader("HashSHA256", string(hashedMetrics)).
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
