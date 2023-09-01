package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oktavarium/go-gauger/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestRootHandle(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name     string
		handlers *Handlers
		url      string
		want     want
	}{
		{
			name:     "simple test on empty path",
			handlers: NewHandlers(storage.NewStorage()),
			url:      "/",
			want: want{
				code:        400,
				response:    "",
				contentType: "",
			},
		},
		{
			name:     "simple test on non empty path",
			handlers: NewHandlers(storage.NewStorage()),
			url:      "/update",
			want: want{
				code:        400,
				response:    "",
				contentType: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, test.url, nil)
			w := httptest.NewRecorder()
			test.handlers.RootHandle(w, r)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			assert.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
		})
	}
}

func TestUpdateHandle(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name     string
		handlers *Handlers
		url      string
		want     want
	}{
		{
			name:     "simple test on wrong metric type",
			handlers: NewHandlers(storage.NewStorage()),
			url:      "/wrongMetricType/",
			want: want{
				code:        400,
				response:    "",
				contentType: "",
			},
		},
		{
			name:     "simple test on empty metric name",
			handlers: NewHandlers(storage.NewStorage()),
			url:      "/gauge/asdf",
			want: want{
				code:        404,
				response:    "",
				contentType: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, test.url, nil)
			w := httptest.NewRecorder()
			test.handlers.UpdateHandle(w, r)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			assert.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
		})
	}
}
