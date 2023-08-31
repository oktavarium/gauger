package server

import (
	"log"
	"net/http"
)

type GaugerServer struct {
	mux  *http.ServeMux
	addr string
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
		// if r.Header.Get("Content-Type") != "text/plain" {
		// 	log.Println("Wrong Content-type in request.")
		// 	w.WriteHeader(http.StatusUnsupportedMediaType)
		// 	return
		// }

		// set right content-type and let's go next
		w.Header().Add("Content-Type", "text/plain")
		next.ServeHTTP(w, r)
	})
}

func NewGaugerServer(addr string) *GaugerServer {
	return &GaugerServer{
		mux:  http.NewServeMux(),
		addr: addr,
	}
}

func (g *GaugerServer) Handle(pattern string, handler http.Handler) {
	g.mux.Handle(pattern, handler)
}

func (g *GaugerServer) ListenAndServe() error {
	return http.ListenAndServe(g.addr, middleware(g.mux))
}
