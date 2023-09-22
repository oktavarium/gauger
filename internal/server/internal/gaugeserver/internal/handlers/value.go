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

	// checking metric name
	if len(metricName) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if metricType == models.GaugeType {
		val, ok := h.storage.GetGauger(metricName)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		valStr := strconv.FormatFloat(val, 'f', -1, 64)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(valStr))
	} else {
		val, ok := h.storage.GetCounter(metricName)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		valStr := strconv.FormatInt(val, 10)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(valStr))
	}
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

	if models.MetricType(metrics.MType) == models.GaugeType {
		val, _ := h.storage.GetGauger(metrics.ID)
		// if !ok {
		// 	w.WriteHeader(http.StatusNotFound)
		// 	return
		// }
		metrics.Value = &val
		encoder := json.NewEncoder(w)
		err := encoder.Encode(&metrics)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		val, _ := h.storage.GetCounter(metrics.ID)
		// if !ok {
		// 	w.WriteHeader(http.StatusNotFound)
		// 	return
		// }
		metrics.Delta = &val
		encoder := json.NewEncoder(w)
		err := encoder.Encode(&metrics)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
