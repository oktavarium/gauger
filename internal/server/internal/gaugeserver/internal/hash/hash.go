package hash

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

	"github.com/oktavarium/go-gauger/internal/server/internal/logger"
	"go.uber.org/zap"
)

type hashedWriter struct {
	w  http.ResponseWriter
	mw io.Writer
	b  *bytes.Buffer
	k  []byte
}

func newHashedWriter(w http.ResponseWriter, key []byte) *hashedWriter {
	b := new(bytes.Buffer)
	mw := io.MultiWriter(w, b)
	return &hashedWriter{w, mw, b, key}
}

func (h *hashedWriter) Header() http.Header {
	return h.w.Header()
}

func (h *hashedWriter) Write(data []byte) (int, error) {
	return h.mw.Write(data)
}

func (h *hashedWriter) WriteHeader(statusCode int) {
	h.w.WriteHeader(statusCode)
}

func (h *hashedWriter) hash() (string, error) {
	return hashData(h.k, h.b.Bytes())
}

func hashData(key []byte, data []byte) (string, error) {
	if len(key) == 0 {
		return hex.EncodeToString(data), nil
	}

	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	hashedData := mac.Sum(nil)
	return hex.EncodeToString(hashedData), nil
}

func checkHash(key []byte, reader io.Reader, hash string) error {
	clientHashBytes, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("error occured on reading from io.Reader: %w", err)
	}
	clientHash, err := hashData(key, clientHashBytes)
	if err != nil {
		return fmt.Errorf("error occured on calculating hash: %w", err)
	}
	if clientHash != hash {
		return fmt.Errorf("hashes are not equal")
	}
	return nil
}

func HashMiddleware(key []byte) func(http.Handler) http.Handler {
	nextF := func(next http.Handler) http.Handler {
		hf := func(w http.ResponseWriter, r *http.Request) {
			clientHash := r.Header.Get("Hashsha256")
			fmt.Println(clientHash, r.Header)
			if len(clientHash) == 0 {
				logger.Logger().Info("error",
					zap.String("func", "HashMiddleware"),
					zap.String("msg", "empty hash"))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			err := checkHash(key, r.Body, clientHash)
			if err != nil {
				logger.Logger().Info("error",
					zap.String("func", "HashMiddleware"),
					zap.Error(err))

				w.WriteHeader(http.StatusBadRequest)
				return
			}

			hashedWriter := newHashedWriter(w, key)

			next.ServeHTTP(hashedWriter, r)

			hash, err := hashedWriter.hash()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			hashedWriter.Header().Set("HashSHA256", hash)
		}
		return http.HandlerFunc(hf)
	}

	return nextF
}
