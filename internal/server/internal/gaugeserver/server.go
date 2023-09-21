package gaugeserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/handlers"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/storage"
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

	server.router.Use(middleware.Logger)
	server.router.Get("/", handler.GetHandle)
	server.router.Post("/update", handler.UpdateHandle)
	server.router.Get("/value", handler.ValueHandle)

	return server
}

func (g *GaugerServer) ListenAndServe() error {
	return http.ListenAndServe(g.addr, g.router)
}
