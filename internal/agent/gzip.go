package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"

	"github.com/oktavarium/go-gauger/internal/shared"
)

func compressMetrics(metrics []shared.Metric) ([]byte, error) {
	var compressedJSON bytes.Buffer
	wr := gzip.NewWriter(&compressedJSON)

	encoder := json.NewEncoder(wr)
	for _, v := range metrics {
		if err := encoder.Encode(v); err != nil {
			return nil, fmt.Errorf("error occured on encoding metric: %w", err)
		}
	}

	wr.Close()
	return compressedJSON.Bytes(), nil
}
