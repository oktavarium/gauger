package handlers

import "net/http"

func (h *Handler) PingHandle(w http.ResponseWriter, r *http.Request) {
	err := h.storage.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
