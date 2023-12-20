package handlers

import (
	"net/http"
	"testing"
)

type resw struct {
	statusCode int
	header     http.Header
	reply      string
}

func (r *resw) Header() http.Header {
	return r.header
}

func (r *resw) Write(reply []byte) (int, error) {
	r.reply = string(reply)
	return len(reply), nil
}

func (r *resw) WriteHeader(statusCode int) {
	r.statusCode = statusCode
}

func newResw() *resw {
	return &resw{
		header: make(http.Header),
	}
}

// GetHandle получает все доступные в данный момент метрики
func TestGetHandle(t *testing.T) {
	// body := ""
	// reader := bytes.NewReader([]byte(body))
	// req, err := http.NewRequest("GET", "localhost", reader)
	// require.NoError(t, err)

	// storage, err := storage.NewInMemoryStorage("/tmp/storage.db", false, 1*time.Minute)
	// require.NoError(t, err)
	// h := NewHandler(storage)

	// rw := newResw()

	// h.GetHandle(rw, req)

	// require.Equal(t, http.StatusOK, rw.statusCode)
}
