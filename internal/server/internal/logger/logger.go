package logger

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

var logger *zap.Logger = zap.NewNop()

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

	}

}
