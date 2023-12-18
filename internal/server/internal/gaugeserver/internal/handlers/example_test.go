package handlers_test

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/oktavarium/go-gauger/internal/shared"
)

// GetHandle получает все доступные в данный момент метрики
func Example() {
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
			return fmt.Errorf("error on hashin data: %w", err)
		}

		request = request.SetHeader("HashSHA256", hash)
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

}
