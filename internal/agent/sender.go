package agent

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/oktavarium/go-gauger/internal/shared"
	"golang.org/x/sync/errgroup"
)

const updatePath string = "updates"

func reportMetrics(address string, key string, metrics []byte) error {
	endpoint := fmt.Sprintf("%s/%s/", address, updatePath)
	var metricsResponse shared.Metric

	client := resty.New()
	request := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip")

	if len(key) != 0 {
		request = request.SetHeader("HashSHA256",
			hashData([]byte(key), metrics))
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
	eg *errgroup.Group,
	d time.Duration,
	inCh chan []byte) {

	eg.Go(func() error {
		for v := range inCh {
			err := reportMetrics(address, key, v)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
