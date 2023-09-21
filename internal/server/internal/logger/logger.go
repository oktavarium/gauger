package logger

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

var logger *zap.Logger = zap.NewNop()

type info struct {
	size   int
	status int
}

type loggedResponseWriter struct {
	w http.ResponseWriter
	i *info
}

func (lrw *loggedResponseWriter) Header() http.Header {
	return lrw.w.Header()
}

func (lrw *loggedResponseWriter) Write(body []byte) (int, error) {
	size, err := lrw.w.Write(body)
	lrw.i.size = size
	return size, err
}

func (lrw *loggedResponseWriter) WriteHeader(statusCode int) {
	lrw.w.WriteHeader(statusCode)
	fmt.Println(statusCode)
	lrw.i.status = statusCode
}

func Logger() *zap.Logger {
	return logger
}

func Init(level string) error {
	atomicLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return fmt.Errorf("error on parsing zap atomic level: %w", err)
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = atomicLevel
	zl, err := cfg.Build()
	if err != nil {
		return fmt.Errorf("error on building zap config: %w", err)
	}

	logger = zl
	return nil
}

func LoggerMiddleware(h http.Handler) http.Handler {
	hf := func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		method := r.Method
		start := time.Now()

		loggerRW := loggedResponseWriter{
			w: w,
			i: &info{},
		}

		h.ServeHTTP(&loggerRW, r)
		duration := time.Since(start)

		Logger().Info(">",
			zap.String("uri", uri),
			zap.String("method", method),
			zap.Int64("duration ms", duration.Milliseconds()),
		)

		Logger().Info("<",
			zap.Int("size", loggerRW.i.size),
			zap.Int("status", loggerRW.i.status),
		)
	}
	return http.HandlerFunc(hf)
}
