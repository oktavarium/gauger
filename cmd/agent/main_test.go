package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMetrics(t *testing.T) {
	tests := []struct {
		name              string
		wantGaugeValues   []string
		wantGaugeType     metricType
		wantCounterValues []string
		wantCounterType   metricType
	}{
		{
			name: "checking gauge types",
			wantGaugeValues: []string{"Alloc",
				"TotalAlloc",
				"Sys",
				"Lookups",
				"Frees",
				"Mallocs",
				"HeapAlloc",
				"HeapSys",
				"HeapIdle",
				"HeapInuse",
				"HeapReleased",
				"HeapObjects",
				"StackInuse",
				"StackSys",
				"MSpanInuse",
				"MSpanSys",
				"MCacheInuse",
				"MCacheSys",
				"BuckHashSys",
				"GCSys",
				"OtherSys",
				"NextGC",
				"LastGC",
				"PauseTotalNs",
				"NumGC",
				"NumForcedGC",
				"GCCPUFraction",
				"RandomValue",
			},
			wantGaugeType:     gaugeType,
			wantCounterValues: []string{"PollCount"},
			wantCounterType:   counterType,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			metrics := NewMetrics()
			assert.Equal(t, test.wantGaugeType, metrics.gauges.mType)
			assert.Equal(t, test.wantCounterType, metrics.counters.mType)
			for _, v := range test.wantGaugeValues {
				_, ok := metrics.gauges.metrics[v]
				assert.Equal(t, true, ok)
			}
			for _, v := range test.wantCounterValues {
				_, ok := metrics.counters.metrics[v]
				assert.Equal(t, true, ok)
			}

		})
	}
}
