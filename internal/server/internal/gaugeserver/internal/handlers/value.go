package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/oktavarium/go-gauger/internal/shared"
)

func (h *Handler) ValueHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	metricType := strings.ToLower(chi.URLParam(r, "type"))
	metricName := strings.ToLower(chi.URLParam(r, "name"))

	var valStr string
	switch metricType {
	case shared.GaugeType:
		val, ok := h.storage.GetGauger(metricName)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		valStr = strconv.FormatFloat(val, 'f', -1, 64)

	case shared.CounterType:
		val, ok := h.storage.GetCounter(metricName)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		valStr = strconv.FormatInt(val, 10)

	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte(valStr))
}

func (h *Handler) ValueJSONHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	var metric shared.Metric
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&metric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// checking metric name
	if len(metric.ID) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch metric.MType {
	case shared.GaugeType:
		val, ok := h.storage.GetGauger(metric.ID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		metric.Value = &val

	case shared.CounterType:
		val, ok := h.storage.GetCounter(metric.ID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		metric.Delta = &val

	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(&metric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
