package agent

import "github.com/oktavarium/go-gauger/internal/models"

type gaugeMetrics struct {
	metrics map[string]float64
	mType   models.MetricType
}

type counterMetrics struct {
	metrics map[string]int64
	mType   models.MetricType
}

type metrics struct {
	gauges   gaugeMetrics
	counters counterMetrics
}
