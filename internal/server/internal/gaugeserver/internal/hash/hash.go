package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

func hashData(key []byte, data []byte) (string, error) {
	if len(key) == 0 {
		return hex.EncodeToString(data), nil
	}

	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	hashedData := mac.Sum(nil)
	return hex.EncodeToString(hashedData), nil
}

func HashMiddleware(key []byte) func(http.Handler) http.Handler {
	hf := func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP()

	}

	return http.HandleFunc(hf)
}
