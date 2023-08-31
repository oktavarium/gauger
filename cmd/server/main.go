package main

import (
	"log"
	"net/http"
	"strings"
)

const srvAddr string = "localhost:8080"

type storage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("New request")
		// check request method
		if r.Method != http.MethodPost {
			log.Println("Wrong request's method.")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// check request content-type
		if r.Header.Get("Content-Type") != "text/plain" {
			log.Println("Wrong Content-type in request.")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		// set right content-type and let's go next
		w.Header().Add("Content-Type", "text/plain")
		next.ServeHTTP(w, r)
	})
}

func rootHandle(w http.ResponseWriter, r *http.Request) {
	log.Println("We don't serve empty path and paths that are not update-paths.")
	w.WriteHeader(http.StatusBadRequest)
	return
}

// parsing URL http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func updateHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	log.Println("We serve update path.")
	urlParams := strings.Split(r.URL.String(), "/")
	log.Println(urlParams)
	return
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()
	mux.Handle(`/`, http.HandlerFunc(rootHandle))
	mux.Handle("/update/", http.HandlerFunc(updateHandle))
	mux.Handle("/update", http.HandlerFunc(updateHandle))

	return http.ListenAndServe(srvAddr, middleware(mux))
}
