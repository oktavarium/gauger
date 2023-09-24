package handlers

import "net/http"

func (h *Handler) GetHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(h.storage.GetAll()))
}
