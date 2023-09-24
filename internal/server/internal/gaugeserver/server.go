package gaugeserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/gzip"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/handlers"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/storage"
	"github.com/oktavarium/go-gauger/internal/server/internal/logger"
)

type GaugerServer struct {
	router *chi.Mux
	addr   string
}

func NewGaugerServer(addr string) *GaugerServer {
	server := &GaugerServer{
		router: chi.NewRouter(),
		addr:   addr,
	}
	storage := storage.NewStorage()
	handler := handlers.NewHandler(storage)

	server.router.Use(logger.LoggerMiddleware)
	server.router.Use(gzip.GzipMiddleware)
	server.router.Get("/", handler.GetHandle)
	server.router.Post("/update/", handler.UpdateJSONHandle)
	server.router.Post("/value/", handler.ValueJSONHandle)
	server.router.Post("/update/{type}/{name}/{value}", handler.UpdateHandle)
	server.router.Get("/value/{type}/{name}", handler.ValueHandle)

	return server
}

func (g *GaugerServer) ListenAndServe() error {
	return http.ListenAndServe(g.addr, g.router)
}
