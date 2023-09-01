package main

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"golang.org/x/exp/rand"
)

const (
	pollInterval   time.Duration = 2 * time.Second
	reportInterval time.Duration = 10 * time.Second
)

type stats struct {
	gaugeTypes   map[string]float64
	counterTypes map[string]int64
}

func NewStats() stats {
	return stats{
		gaugeTypes: map[string]float64{
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
		counterTypes: map[string]int64{
			"PollCount": 0,
		},
	}
}

func main() {
	fmt.Println("Agent started")
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	stats := NewStats()
	for {
		for i := 0; i < 5; i++ {
			statsReader(&stats)
			time.Sleep(pollInterval)
		}
		if err := reportStats(&stats); err != nil {
			return err
		}
	}
}

func statsReader(st *stats) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	st.gaugeTypes["Alloc"] = float64(memStats.Alloc)
	st.gaugeTypes["TotalAlloc"] = float64(memStats.TotalAlloc)
	st.gaugeTypes["Sys"] = float64(memStats.Sys)
	st.gaugeTypes["Lookups"] = float64(memStats.Lookups)
	st.gaugeTypes["Frees"] = float64(memStats.Frees)
	st.gaugeTypes["Mallocs"] = float64(memStats.Mallocs)
	st.gaugeTypes["HeapAlloc"] = float64(memStats.HeapAlloc)
	st.gaugeTypes["HeapSys"] = float64(memStats.HeapSys)
	st.gaugeTypes["HeapIdle"] = float64(memStats.HeapIdle)
	st.gaugeTypes["HeapInuse"] = float64(memStats.HeapInuse)
	st.gaugeTypes["HeapReleased"] = float64(memStats.HeapReleased)
	st.gaugeTypes["HeapObjects"] = float64(memStats.HeapObjects)
	st.gaugeTypes["StackInuse"] = float64(memStats.StackInuse)
	st.gaugeTypes["StackSys"] = float64(memStats.StackSys)
	st.gaugeTypes["MSpanInuse"] = float64(memStats.MSpanInuse)
	st.gaugeTypes["MSpanSys"] = float64(memStats.MSpanSys)
	st.gaugeTypes["MCacheInuse"] = float64(memStats.MCacheInuse)
	st.gaugeTypes["MCacheSys"] = float64(memStats.MCacheSys)
	st.gaugeTypes["BuckHashSys"] = float64(memStats.BuckHashSys)
	st.gaugeTypes["GCSys"] = float64(memStats.GCSys)
	st.gaugeTypes["OtherSys"] = float64(memStats.OtherSys)
	st.gaugeTypes["NextGC"] = float64(memStats.NextGC)
	st.gaugeTypes["LastGC"] = float64(memStats.LastGC)
	st.gaugeTypes["PauseTotalNs"] = float64(memStats.PauseTotalNs)
	st.gaugeTypes["NumGC"] = float64(memStats.NumGC)
	st.gaugeTypes["NumForcedGC"] = float64(memStats.NumForcedGC)
	st.gaugeTypes["GCCPUFraction"] = float64(memStats.GCCPUFraction)

	st.gaugeTypes["RandomValue"] = rand.Float64()
	st.counterTypes["PollCount"]++
}

func reportStats(st *stats) error {
	client := &http.Client{}
	for k, v := range st.gaugeTypes {
		path := fmt.Sprintf("http://localhost:8080/update/%s/%s/%f", "gauge", k, v)
		if err := makeRequest(client, path); err != nil {
			return err
		}
	}

	for k, v := range st.counterTypes {
		path := fmt.Sprintf("http://localhost:8080/update/%s/%s/%d", "counter", k, v)
		if err := makeRequest(client, path); err != nil {
			return err
		}
	}

	return nil
}

func makeRequest(client *http.Client, path string) error {
	req, err := http.NewRequest("POST", path, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("Response is not OK")
	}
	return nil
}
