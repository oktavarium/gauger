package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oktavarium/go-gauger/internal/storage"
	"github.com/stretchr/testify/assert"
)

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
			name:     "check wrong metric type",
			handlers: NewHandlers(storage.NewStorage()),
			url:      "/update/wrongMetricType/test/0.0",
			want: want{
				code:        400,
				response:    "",
				contentType: "",
			},
		},
		{
			name:     "check wrong metric name",
			handlers: NewHandlers(storage.NewStorage()),
			url:      "/update/gauge/",
			want: want{
				code:        404,
				response:    "",
				contentType: "",
			},
		},
		{
			name:     "check wrong metric gauge value",
			handlers: NewHandlers(storage.NewStorage()),
			url:      "/update/gauge/Alloc/WrongVal",
			want: want{
				code:        400,
				response:    "",
				contentType: "",
			},
		},
		{
			name:     "check wrong metric counter value",
			handlers: NewHandlers(storage.NewStorage()),
			url:      "/update/counter/PollCount/wrongVal",
			want: want{
				code:        400,
				response:    "",
				contentType: "",
			},
		},
		{
			name:     "check right metric counter value",
			handlers: NewHandlers(storage.NewStorage()),
			url:      "/update/counter/PollCount/0",
			want: want{
				code:        200,
				response:    "",
				contentType: "",
			},
		},
		{
			name:     "check right metric gauge value",
			handlers: NewHandlers(storage.NewStorage()),
			url:      "/update/gauge/Alloc/0.0",
			want: want{
				code:        200,
				response:    "",
				contentType: "",
			},
		},
		{
			name:     "check wrong update method",
			handlers: NewHandlers(storage.NewStorage()),
			url:      "/wrongMethod/gauge/Alloc/0.0",
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
