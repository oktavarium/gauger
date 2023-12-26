package handlers

import (
	"net/http"

	"github.com/oktavarium/go-gauger/internal/server/internal/logger"
)

// Ping - определяет доступность хранилища (и самого сервиса)
func (h *Handler) PingHandle(w http.ResponseWriter, r *http.Request) {
	var err error
	err = h.storage.Ping(r.Context())
	if err != nil {
		if err != nil {
			logger.LogError("PingHandle", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
