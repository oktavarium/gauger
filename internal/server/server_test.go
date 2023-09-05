package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oktavarium/go-gauger/internal/handlers"
	"github.com/oktavarium/go-gauger/internal/storage"
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
	server := NewGaugerServer("localhost", handlers.NewHandler(storage.NewStorage()))
	ts := httptest.NewServer(server.router)
	defer ts.Close()

	tests := []struct {
		method string
		url    string
		want   string
		status int
	}{
		{"POST", "/update/gauge/alloc/4.0", "", 200},
		{"POST", "/wrongMethod/gauge/alloc/4.0", "404 page not found\n", 404},
		{"POST", "/update/counter/pollscounter/1", "", 200},
		{"POST", "/update/counter/pollscounter/f", "", 400},
		{"GET", "/value/gauge/alloc", "4", 200},
		{"GET", "/value/counter/pollscounter", "1", 200},
		{"GET", "/value/counter/wrong", "", 404},
		{"GET", "/value/wrongtype/wrong", "", 400},
		{"POST", "/value/counter/pollscounter", "", 405},
		{"GET", "/", "alloc: 4\npollscounter: 1\n", 200},
	}

	for _, test := range tests {
		resp, get := testRequest(t, ts, test.method, test.url)
		defer resp.Body.Close()
		assert.Equal(t, test.status, resp.StatusCode)
		assert.Equal(t, test.want, get)
	}
}
