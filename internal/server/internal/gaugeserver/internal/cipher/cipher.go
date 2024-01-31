package cipher

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/oktavarium/go-gauger/internal/server/internal/logger"
	"go.uber.org/zap"
)

type Cipher struct {
	pk *rsa.PrivateKey
}

func NewCipher(pkFile string) (Cipher, error) {
	cipher := Cipher{}
	privateKeyData, err := os.ReadFile(pkFile)
	if err != nil {
		return cipher, fmt.Errorf("error on reading private key file: %w", err)
	}

	pkPEM, _ := pem.Decode(privateKeyData)
	if pkPEM.Type != "RSA PRIVATE KEY" {
		return cipher, fmt.Errorf("wrong key type: %w", err)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(pkPEM.Bytes)
	if err != nil {
		return cipher, fmt.Errorf("error parsing private key: %w", err)
	}

	cipher.pk = privateKey

	return cipher, nil
}

// CipherMiddleware - метод посредника для деширофки данных
func (c Cipher) CipherMiddleware(next http.Handler) http.Handler {
	hf := func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/octet-stream" {
			next.ServeHTTP(w, r)
			return
		}

		encryptedData, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Logger().Error("error",
				zap.String("func", "CipherMiddleware"),
				zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		opts := rsa.OAEPOptions{
			Hash:    crypto.SHA256,
			MGFHash: 0,
			Label:   []byte{},
		}
		decryptedData, err := c.pk.Decrypt(rand.Reader, encryptedData, opts)
		if err != nil {
			logger.Logger().Error("error",
				zap.String("func", "CipherMiddleware"),
				zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		reader := io.NopCloser(bytes.NewReader(decryptedData))
		r.Body = reader
		defer func() {
			if err = reader.Close(); err != nil {
				logger.Logger().Error("error",
					zap.String("func", "CipherMiddleware"),
					zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(hf)
}
