package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oktavarium/go-gauger/internal/handlers"
)

type GaugerServer struct {
	router *chi.Mux
	addr   string
}

func NewGaugerServer(addr string, handler *handlers.Handler) *GaugerServer {
	server := &GaugerServer{
		router: chi.NewRouter(),
		addr:   addr,
	}
	server.router.Get("/", handler.GetHandle)
	server.router.Post("/update/{type}/{name}/{value}", handler.UpdateHandle)
	server.router.Get("/value/{type}/{name}", handler.ValueHandle)

	return server
}

func (g *GaugerServer) ListenAndServe() error {
	return http.ListenAndServe(g.addr, g.router)
}
