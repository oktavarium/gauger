package handlers

import (
	"net/http"
	"strconv"
	"strings"
)

type metricType string

// declare interfaces where we plan to use it
type metricsSaver interface {
	SaveGauge(namse string, val float64)
	UpdateCounter(name string, val int64)
	CheckMetricName(name string) bool
}

const (
	gaugeType   metricType = "gauge"
	counterType metricType = "counter"
	handleType  string     = "update"
)

type Handlers struct {
	storage metricsSaver
}

func NewHandlers(storage metricsSaver) *Handlers {
	return &Handlers{
		storage: storage,
	}
}

// parsing URL http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func (h *Handlers) UpdateHandle(w http.ResponseWriter, r *http.Request) {
	urlParams := strings.FieldsFunc(r.URL.RequestURI(), func(c rune) bool {
		return c == '/'
	})

	// wrong number of parameters
	for len(urlParams) < 4 {
		urlParams = append(urlParams, "")
	}

	hType := urlParams[0]
	metricType := metricType(urlParams[1])
	metricName := urlParams[2]
	metricValueStr := urlParams[3]

	// checking update method
	if hType != handleType {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// checking metric type
	if metricType != gaugeType && metricType != counterType {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// checking metric name
	if len(metricName) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// checking metric value
	if len(metricValueStr) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if metricType == gaugeType {
		val, err := strconv.ParseFloat(metricValueStr, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.storage.SaveGauge(metricName, val)

	} else {
		val, err := strconv.ParseInt(metricValueStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.storage.UpdateCounter(metricName, val)
	}

	w.WriteHeader(http.StatusOK)
}
