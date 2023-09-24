package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"

	"github.com/oktavarium/go-gauger/internal/models"
)

func compressMetrics(metrics models.Metrics) ([]byte, error) {
	jsonRaw, err := json.Marshal(metrics)
	if err != nil {
		return nil, fmt.Errorf("error on marshaling metrics: %w", err)
	}
	var compressedJSON bytes.Buffer
	wr, err := gzip.NewWriterLevel(&compressedJSON, gzip.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("error on creating gzip writer: %w", err)
	}

	_, err = wr.Write(jsonRaw)
	if err != nil {
		return nil, fmt.Errorf("error on compressing data: %w", err)
	}

	return compressedJson.Bytes(), nil
}
