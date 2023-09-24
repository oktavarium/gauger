package gaugeserver

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/archivarius"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/gzip"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/handlers"
	"github.com/oktavarium/go-gauger/internal/server/internal/logger"
)

type GaugerServer struct {
	router *chi.Mux
	addr   string
}

func NewGaugerServer(addr string,
	filename string,
	restore bool,
	timeout int) (*GaugerServer, error) {
	server := &GaugerServer{
		router: chi.NewRouter(),
		addr:   addr,
	}

	archiver, err := archivarius.NewArchiver(filename, restore, timeout)
	if err != nil {
		return nil, fmt.Errorf("error on creating archiver: %w", err)
	}

	handler := handlers.NewHandler(archiver)

	server.router.Use(logger.LoggerMiddleware)
	server.router.Use(gzip.GzipMiddleware)
	server.router.Get("/", handler.GetHandle)
	server.router.Post("/update/", handler.UpdateJSONHandle)
	server.router.Post("/value/", handler.ValueJSONHandle)
	server.router.Post("/update/{type}/{name}/{value}", handler.UpdateHandle)
	server.router.Get("/value/{type}/{name}", handler.ValueHandle)

	return server, nil
}

func (g *GaugerServer) ListenAndServe() error {
	return http.ListenAndServe(g.addr, g.router)
}
