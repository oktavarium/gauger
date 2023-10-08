package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/oktavarium/go-gauger/internal/shared"
)

func (h *Handler) UpdatesHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	metricType := strings.ToLower(chi.URLParam(r, "type"))
	metricName := strings.ToLower(chi.URLParam(r, "name"))
	metricValueStr := chi.URLParam(r, "value")

	var err error
	switch metricType {
	case shared.GaugeType:
		var val float64
		val, err = strconv.ParseFloat(metricValueStr, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = h.storage.SaveGauge(metricName, val)

	case shared.CounterType:
		var val int64
		val, err = strconv.ParseInt(metricValueStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = h.storage.UpdateCounter(r.Context(), metricName, val)

	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdatesJSONHandle(w http.ResponseWriter, r *http.Request) {
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

	var delta int64
	switch metric.MType {
	case shared.GaugeType:
		err = h.storage.SaveGauge(metric.ID, *metric.Value)

	case shared.CounterType:
		delta, err = h.storage.UpdateCounter(r.Context(), metric.ID, *metric.Delta)
		metric.Delta = &delta

	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(&metric)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
