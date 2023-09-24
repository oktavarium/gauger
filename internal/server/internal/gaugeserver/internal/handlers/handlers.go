package handlers

import (
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/archivarius"
)

type Handler struct {
	archiver archivarius.Archivarius
}

func NewHandler(archiver archivarius.Archivarius) *Handler {
	return &Handler{
		archiver: archiver,
	}
}
