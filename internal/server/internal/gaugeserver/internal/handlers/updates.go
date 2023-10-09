package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/oktavarium/go-gauger/internal/shared"
)

func (h *Handler) UpdatesHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	var metrics []shared.Metric
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&metrics); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(metrics) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.storage.BatchUpdate(r.Context(), metrics)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	metricStab := shared.Metric{}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(metricStab)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
