package handlers

import (
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/storage"
)

type metricType string

const (
	requestMethod string     = "POST"
	gaugeType     metricType = "gauge"
	counterType   metricType = "counter"
	updatePath    string     = "update"
)

type Handler struct {
	storage storage.Storage
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}
