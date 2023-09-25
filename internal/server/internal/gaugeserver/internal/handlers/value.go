package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/oktavarium/go-gauger/internal/models"
)

func (h *Handler) ValueHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	metricType := models.MetricType(strings.ToLower(chi.URLParam(r, "type")))
	metricName := strings.ToLower(chi.URLParam(r, "name"))

	// checking metric type
	if metricType != models.GaugeType && metricType != models.CounterType {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var valStr string
	switch metricType {
	case models.GaugeType:
		val, ok := h.archiver.GetGauger(metricName)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		valStr = strconv.FormatFloat(val, 'f', -1, 64)

	case models.CounterType:
		val, ok := h.archiver.GetCounter(metricName)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		valStr = strconv.FormatInt(val, 10)
	}

	w.Write([]byte(valStr))
}

func (h *Handler) ValueJSONHandle(w http.ResponseWriter, r *http.Request) {
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

	switch models.MetricType(metrics.MType) {
	case models.GaugeType:
		val, ok := h.archiver.GetGauger(metrics.ID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		metrics.Value = &val

	case models.CounterType:
		val, ok := h.archiver.GetCounter(metrics.ID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		metrics.Delta = &val
	}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(&metrics)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
