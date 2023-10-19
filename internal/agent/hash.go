package agent

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
