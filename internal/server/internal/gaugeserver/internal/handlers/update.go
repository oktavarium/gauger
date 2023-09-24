package handlers

import (
	"encoding/json"
	"fmt"
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

	if metricType == models.GaugeType {
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
		fmt.Println(err)
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

	if models.MetricType(metrics.MType) == models.GaugeType {
		h.storage.SaveGauge(metrics.ID, *metrics.Value)
		encoder := json.NewEncoder(w)
		err := encoder.Encode(&metrics)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	} else {
		h.storage.UpdateCounter(metrics.ID, *metrics.Delta)
		val, _ := h.storage.GetCounter(metrics.ID)
		metrics.Delta = &val
		encoder := json.NewEncoder(w)
		err := encoder.Encode(&metrics)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
