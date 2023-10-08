package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/oktavarium/go-gauger/internal/shared"
)

func (h *Handler) UpdatesHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	metrics := make([]shared.Metric, 0)
	var metric shared.Metric
	decoder := json.NewDecoder(r.Body)
	for {
		if err := decoder.Decode(&metric); err == io.EOF {
			break
		} else if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metrics = append(metrics, metric)
	}

	err := h.storage.BatchUpdate(r.Context(), metrics)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
