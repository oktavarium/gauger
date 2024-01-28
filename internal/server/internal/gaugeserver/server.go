package gaugeserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/cipher"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/gzip"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/handlers"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/hash"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/storage"
	"github.com/oktavarium/go-gauger/internal/server/internal/logger"
	"go.uber.org/zap"
)

// GaugeServer - управляющий сервис для сбора метрик
type GaugerServer struct {
	ctx    context.Context
	router *chi.Mux
	srv    *http.Server
}

// NewGaugeServer - конструктор сервиса ядл сбора метрик
func NewGaugerServer(
	ctx context.Context,
	addr string,
	filename string,
	restore bool,
	timeout time.Duration,
	dsn string,
	key string,
	pkFile string) (*GaugerServer, error) {
	server := &GaugerServer{
		ctx:    ctx,
		router: chi.NewRouter(),
	}

	var s storage.Storage
	var err error

	if len(dsn) == 0 {
		s, err = storage.NewInMemoryStorage(ctx, filename, restore, timeout)
	} else {
		s, err = storage.NewPostgresqlStorage(dsn)
	}
	if err != nil {
		return nil, fmt.Errorf("error on creating storage: %w", err)
	}

	handler := handlers.NewHandler(s)

	if len(pkFile) != 0 {
		c, err := cipher.NewCipher(pkFile)
		if err != nil {
			return nil, fmt.Errorf("error on creating cipher: %w", err)
		}
		server.router.Use(c.CipherMiddleware)
	}

	server.router.Use(logger.LoggerMiddleware)
	if len(key) != 0 {
		server.router.Use(hash.HashMiddleware([]byte(key)))
	}

	server.router.Use(gzip.GzipMiddleware)
	server.router.Get("/", handler.GetHandle)
	server.router.Get("/ping", handler.PingHandle)
	server.router.Post("/update/", handler.UpdateJSONHandle)
	server.router.Post("/updates/", handler.UpdatesHandle)
	server.router.Post("/value/", handler.ValueJSONHandle)
	server.router.Post("/update/{type}/{name}/{value}", handler.UpdateHandle)
	server.router.Get("/value/{type}/{name}", handler.ValueHandle)

	server.srv = &http.Server{
		Addr:    addr,
		Handler: server.router,
	}

	return server, nil
}

// ListenAndServer - запуск сервиса
func (g *GaugerServer) ListenAndServe() error {
	go func() {
		<-g.ctx.Done()
		if err := g.srv.Shutdown(context.Background()); err != nil {
			logger.Logger().Error("error",
				zap.String("func", "ListenAndServer"),
				zap.Error(err))
		}

	}()

	if err := g.srv.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("error on listen and serve: %w", err)
	}

	return nil
}
