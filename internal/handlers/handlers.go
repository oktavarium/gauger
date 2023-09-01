package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/oktavarium/go-gauger/internal/storage"
)

type metricType string

const (
	gaugeType   metricType = "gauge"
	counterType metricType = "counter"
)

type Handlers struct {
	storage storage.MetricsSaver
}

func NewHandlers(storage storage.MetricsSaver) *Handlers {
	return &Handlers{
		storage: storage,
	}
}

func (h *Handlers) RootHandle(w http.ResponseWriter, r *http.Request) {
	log.Println("We don't serve empty path and paths that are not update-paths.")
	w.WriteHeader(http.StatusBadRequest)
}

// parsing URL http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func (h *Handlers) UpdateHandle(w http.ResponseWriter, r *http.Request) {
	log.Println("We serve update path.")
	urlParams := strings.Split(r.URL.RequestURI(), "/")
	for len(urlParams) < 3 {
		urlParams = append(urlParams, "")
	}

	metricType := metricType(urlParams[0])
	metricName := urlParams[1]
	metricValueStr := urlParams[2]
	fmt.Println(urlParams, "!", urlParams[0], "!", len(urlParams), r.URL.RequestURI())
	if metricType != gaugeType && metricType != counterType {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(metricName) == 0 {
		w.WriteHeader(http.StatusNotFound)
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
