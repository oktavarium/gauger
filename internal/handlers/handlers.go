package handlers

import (
	"github.com/oktavarium/go-gauger/internal/storage"
)

type metricType string

const (
	gaugeType   metricType = "gauge"
	counterType metricType = "counter"
	handleType  string     = "update"
)

type Handler struct {
	storage storage.Storage
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}
