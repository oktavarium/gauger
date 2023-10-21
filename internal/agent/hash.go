package agent

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func hashData(key []byte, data []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	hashedData := mac.Sum(nil)
	return hex.EncodeToString(hashedData)
}
