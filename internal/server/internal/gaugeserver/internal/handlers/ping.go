package handlers

import (
	"net/http"

	"github.com/oktavarium/go-gauger/internal/server/internal/logger"
	"go.uber.org/zap"
)

// Ping - определяет доступность хранилища (и самого сервиса)
func (h *Handler) PingHandle(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			logger.Logger().Info("error",
				zap.String("func", "PingHandle"),
				zap.Error(err),
			)
		}
	}()
	err = h.storage.Ping(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
