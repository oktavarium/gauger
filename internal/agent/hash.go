package agent

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func hashData(key []byte, data []byte) (string, error) {
	mac := hmac.New(sha256.New, key)
	if _, err := mac.Write(data); err != nil {
		return "",
			fmt.Errorf("error on writing data to hash writer: %w", err)
	}

	hashedData := mac.Sum(nil)
	return hex.EncodeToString(hashedData), nil
}
