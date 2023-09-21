package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/oktavarium/go-gauger/internal/models"
)

func (h *Handler) UpdateHandle(w http.ResponseWriter, r *http.Request) {
	if w.Header().Get("Content-Type") != "application/json" {
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
	if metricType(metrics.MType) != gaugeType && metricType(metrics.MType) != counterType {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// checking metric name
	if len(metrics.ID) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if metricType(metrics.MType) == gaugeType {
		h.storage.SaveGauge(metrics.ID, *metrics.Value)
		encoder := json.NewEncoder(w)
		err := encoder.Encode(&metrics)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	} else {
		h.storage.UpdateCounter(metrics.ID, *metrics.Delta)
		encoder := json.NewEncoder(w)
		err := encoder.Encode(&metrics)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
