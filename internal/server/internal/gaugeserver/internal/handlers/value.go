package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) ValueHandle(w http.ResponseWriter, r *http.Request) {
	metricType := metricType(strings.ToLower(chi.URLParam(r, "type")))
	metricName := strings.ToLower(chi.URLParam(r, "name"))

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

	if metricType == gaugeType {
		val, ok := h.storage.GetGauger(metricName)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		valStr := strconv.FormatFloat(val, 'f', -1, 64)
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
