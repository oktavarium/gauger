package api

import (
	"github.com/oktavarium/go-gauger/internal/shared"
)

func ConvertDBMetricsToMetrics(metrics []shared.Metric) []*Metric {
	allMetrics := make([]*Metric, 0, len(metrics))
	for _, m := range metrics {
		allMetrics = append(allMetrics, ConvertDBMetricToMetric(m))
	}

	return allMetrics
}

func ConvertDBMetricToMetric(metric shared.Metric) *Metric {
	m := &Metric{
		Id:   metric.ID,
		Type: metric.MType,
	}

	switch metric.MType {
	case shared.GaugeType:
		m.Value = float64(*metric.Value)
	case shared.CounterType:
		m.Value = float64(*metric.Delta)
	}

	return m
}

func ConvertMetricsToDBMetrics(metrics []*Metric) []shared.Metric {
	allMetrics := make([]shared.Metric, 0, len(metrics))
	for _, m := range metrics {
		allMetrics = append(allMetrics, ConvertMetricToDBMetric(m))
	}

	return allMetrics
}

func ConvertMetricToDBMetric(metric *Metric) shared.Metric {
	m := shared.Metric{
		ID:    metric.GetId(),
		MType: metric.GetType(),
	}

	switch metric.GetType() {
	case shared.GaugeType:
		m.Value = &metric.Value
	case shared.CounterType:
		v := int64(metric.Value)
		m.Delta = &v
	}

	return m
}
