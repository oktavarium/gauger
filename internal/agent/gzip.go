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
	wr := gzip.NewWriter(&compressedJSON)

	_, err = wr.Write(jsonRaw)
	if err != nil {
		return nil, fmt.Errorf("error on compressing data: %w", err)
	}
	wr.Close()
	return compressedJSON.Bytes(), nil
}
