package shared

type MetricType string

const (
	GaugeType   string = "gauge"
	CounterType string = "counter"
)

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func NewGaugeMetric(id string, val *float64) Metric {
	return Metric{
		ID:    id,
		MType: GaugeType,
		Value: val,
	}
}

func NewCounterMetric(id string, val *int64) Metric {
	return Metric{
		ID:    id,
		MType: CounterType,
		Delta: val,
	}
}

func NewEmptyGaugeMetric() Metric {
	return Metric{
		MType: GaugeType,
	}
}

func NewEmptyCounterMetric() Metric {
	return Metric{
		MType: CounterType,
	}
}
