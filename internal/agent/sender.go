package agent

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/oktavarium/go-gauger/internal/agent/internal/logger"
	"github.com/oktavarium/go-gauger/internal/shared"
)

const updatePath string = "updates"

func reportMetrics(address string, key string, pk *rsa.PublicKey, metrics []byte) error {
	var err error
	endpoint := fmt.Sprintf("%s/%s/", address, updatePath)
	var metricsResponse shared.Metric
	client := resty.New()
	request := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip")

	if len(key) != 0 {
		hash, err := hashData([]byte(key), metrics)
		if err != nil {
			return fmt.Errorf("error on hashing data: %w", err)
		}

		request = request.SetHeader("HashSHA256", hash)
	}

	if pk != nil {
		metrics, err = rsa.EncryptOAEP(sha256.New(), rand.Reader, pk, metrics, []byte{})
		if err != nil {
			return fmt.Errorf("error on encrypting metrics: %w", err)
		}
	}

	request = request.
		SetBody(metrics).
		SetResult(&metricsResponse)
	resp, err := request.Post(endpoint)

	if err != nil {
		return fmt.Errorf("error on making update request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.New("response status code is not OK (200)")
	}

	return nil
}

func sender(ctx context.Context,
	address string,
	key string,
	pk *rsa.PublicKey,
	d time.Duration,
	inCh <-chan []byte) error {

	for {
		select {
		case v := <-inCh:
			if err := reportMetrics(address, key, pk, v); err != nil {
				logger.LogError("error on reporting metrics", err)
			}

		case <-ctx.Done():
			return nil
		}
	}
}
