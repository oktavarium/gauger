package main

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"golang.org/x/exp/rand"
)

type metricType string

const (
	requestMethod string     = "POST"
	gaugeType     metricType = "gauge"
	counterType   metricType = "counter"
	updatePath    string     = "update"
)

type gaugeMetrics struct {
	metrics map[string]float64
	mType   metricType
}

type counterMetrics struct {
	metrics map[string]int64
	mType   metricType
}

type metrics struct {
	gauges   gaugeMetrics
	counters counterMetrics
}

func NewMetrics() metrics {
	return metrics{
		gauges: gaugeMetrics{
			metrics: map[string]float64{
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
			mType: gaugeType,
		},
		counters: counterMetrics{
			metrics: map[string]int64{
				"PollCount": 0,
			},
			mType: counterType,
		},
	}
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	err := parseFlags()
	if err != nil {
		return err
	}

	metrics := NewMetrics()
	var sleepCounter int
	for {
		time.Sleep(1 * time.Second)
		sleepCounter++
		if sleepCounter%flagPollInterval == 0 {
			statsReader(&metrics)
		}
		if sleepCounter%flagReportInterval == 0 {
			if err := reportMetrics(&metrics); err != nil {
				panic(err)
			}
		}
	}
}

func statsReader(m *metrics) {
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
	m.gauges.mType = gaugeType

	m.counters.metrics["PollCount"]++
	m.counters.mType = counterType
}

func reportMetrics(m *metrics) error {
	for k, v := range m.gauges.metrics {
		err := makeUpdateRequest(fmt.Sprintf("%s/%s/%s/%f", flagEndpointAddr+updatePath,
			string(m.gauges.mType), k, v))
		if err != nil {
			return err
		}
	}

	for k, v := range m.counters.metrics {
		err := makeUpdateRequest(fmt.Sprintf("%s/%s/%s/%d", flagEndpointAddr+updatePath,
			string(m.counters.mType), k, v))
		if err != nil {
			return err
		}
	}

	return nil
}

func makeUpdateRequest(endpoint string) error {
	resp, err := http.Post(endpoint, "", nil)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("response is not OK")
	}
	return nil
}
