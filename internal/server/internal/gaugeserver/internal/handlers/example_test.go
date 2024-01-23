package handlers_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/handlers"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/storage"
)

func Example() {
	s, _ := storage.NewInMemoryStorage("/tmp/testdb", false, 1*time.Minute)
	h := handlers.NewHandler(s)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/update/counter/myCounter/2", nil)

	h.UpdateHandle(w, r)

	r = httptest.NewRequest(http.MethodGet, "/value/counter/myCounter", nil)

	h.UpdateHandle(w, r)

	val, _ := io.ReadAll(w.Body)
	fmt.Print(string(val))
}
