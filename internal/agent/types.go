package agent

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
