package gaugeserver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server,
	method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	server, _ := NewGaugerServer("localhost", "/tmp/tmp.txt", false, 0, "", "")
	ts := httptest.NewServer(server.router)
	defer ts.Close()

	tests := []struct {
		name   string
		method string
		url    string
		want   string
		status int
	}{
		{"simple test 1", "POST", "/update/gauge/alloc/4.0", "", 200},
		{"simple test 2", "POST", "/wrongMethod/gauge/alloc/4.0", "404 page not found\n", 404},
		{"simple test 3", "POST", "/update/counter/pollscounter/1", "", 200},
		{"simple test 4", "POST", "/update/counter/pollscounter/f", "", 400},
		{"simple test 5", "GET", "/value/gauge/alloc", "4", 200},
		{"simple test 6", "GET", "/value/counter/pollscounter", "1", 200},
		{"simple test 7", "GET", "/value/counter/wrong", "", 404},
		{"simple test 8", "GET", "/value/wrongtype/wrong", "", 400},
		{"simple test 9", "POST", "/value/counter/pollscounter", "", 405},
		{"simple test 10", "GET", "/", "{\"id\":\"alloc\",\"type\":\"gauge\",\"value\":4}\n{\"id\":\"pollscounter\",\"type\":\"counter\",\"delta\":1}\n", 200},
	}

	for _, test := range tests {
		resp, get := testRequest(t, ts, test.method, test.url)
		defer resp.Body.Close()
		assert.Equal(t, test.status, resp.StatusCode, test.name)
		assert.Equal(t, test.want, get, test.name)
	}
}

func TestNewGaugeServer(t *testing.T) {
	_, err := NewGaugerServer(":8080", "tmp.file", false, 1*time.Second, "", "key")

	require.NoError(t, err)
}
