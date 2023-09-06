package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) UpdateHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	metricType := metricType(strings.ToLower(chi.URLParam(r, "type")))
	metricName := strings.ToLower(chi.URLParam(r, "name"))
	metricValueStr := chi.URLParam(r, "value")

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
