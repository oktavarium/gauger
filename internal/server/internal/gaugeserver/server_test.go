package gaugeserver

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewGaugeServer(t *testing.T) {
	_, err := NewGaugerServer(":8080", "tmp.file", false, 1*time.Second, "", "key")

	require.NoError(t, err)
}
