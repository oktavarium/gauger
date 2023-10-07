package handlers

import "net/http"

func (h *Handler) PingHandle(w http.ResponseWriter, r *http.Request) {
	err := h.storage.Ping(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
