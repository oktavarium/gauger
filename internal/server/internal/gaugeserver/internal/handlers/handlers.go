package handlers

import (
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/storage"
)

type Handler struct {
	storage storage.Storage
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}
