package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/oktavarium/go-gauger/internal/models"
)

func (h *Handler) UpdateHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	metricType := models.MetricType(strings.ToLower(chi.URLParam(r, "type")))
	metricName := strings.ToLower(chi.URLParam(r, "name"))
	metricValueStr := chi.URLParam(r, "value")

	// checking metric type
	if metricType != models.GaugeType && metricType != models.CounterType {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var err error
	switch metricType {
	case models.GaugeType:
		var val float64
		val, err = strconv.ParseFloat(metricValueStr, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = h.archiver.SaveGauge(metricName, val)

	case models.CounterType:
		var val int64
		val, err = strconv.ParseInt(metricValueStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = h.archiver.UpdateCounter(metricName, val)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateJSONHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	var metrics models.Metrics
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&metrics)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// checking metric type
	if models.MetricType(metrics.MType) != models.GaugeType && models.MetricType(metrics.MType) != models.CounterType {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// checking metric name
	if len(metrics.ID) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var delta int64
	switch models.MetricType(metrics.MType) {
	case models.GaugeType:
		err = h.archiver.SaveGauge(metrics.ID, *metrics.Value)
	case models.CounterType:
		delta, err = h.archiver.UpdateCounter(metrics.ID, *metrics.Delta)
		metrics.Delta = &delta
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(&metrics)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
