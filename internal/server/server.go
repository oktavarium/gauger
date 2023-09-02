package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type metricType string

// declare interfaces where we plan to use it
type metricsSaver interface {
	SaveGauge(namse string, val float64)
	UpdateCounter(name string, val int64)
	GetGauger(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
	GetAll() string
}

const (
	gaugeType   metricType = "gauge"
	counterType metricType = "counter"
	handleType  string     = "update"
)

type GaugerServer struct {
	router  *chi.Mux
	addr    string
	storage metricsSaver
}

func NewGaugerServer(addr string, storage metricsSaver) *GaugerServer {
	server := &GaugerServer{
		router:  chi.NewRouter(),
		addr:    addr,
		storage: storage,
	}
	server.router.Get("/", server.getHandle)
	server.router.Post("/update/{type}/{name}/{value}", server.updateHandle)
	server.router.Get("/value/{type}/{name}/", server.valueHandle)

	return server
}

func (g *GaugerServer) ListenAndServe() error {
	return http.ListenAndServe(g.addr, g.router)
}

func (g *GaugerServer) getHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(g.storage.GetAll()))
	w.WriteHeader(http.StatusOK)
}

func (g *GaugerServer) updateHandle(w http.ResponseWriter, r *http.Request) {
	metricType := metricType(strings.ToLower(chi.URLParam(r, "type")))
	metricName := strings.ToLower(chi.URLParam(r, "name"))
	metricValueStr := chi.URLParam(r, "value")

	// checking metric type
	if metricType != gaugeType && metricType != counterType {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// checking metric name
	if len(metricName) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// checking metric value
	if len(metricValueStr) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if metricType == gaugeType {
		val, err := strconv.ParseFloat(metricValueStr, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		g.storage.SaveGauge(metricName, val)

	} else {
		val, err := strconv.ParseInt(metricValueStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		g.storage.UpdateCounter(metricName, val)
	}

	w.WriteHeader(http.StatusOK)
}

func (g *GaugerServer) valueHandle(w http.ResponseWriter, r *http.Request) {
	metricType := metricType(strings.ToLower(chi.URLParam(r, "type")))
	metricName := strings.ToLower(chi.URLParam(r, "name"))

	// checking metric type
	if metricType != gaugeType && metricType != counterType {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// checking metric name
	if len(metricName) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if metricType == gaugeType {
		val, ok := g.storage.GetGauger(metricName)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		valStr := strconv.FormatFloat(val, 'f', -1, 64)
		w.Write([]byte(valStr))
	} else {
		val, ok := g.storage.GetCounter(metricName)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		valStr := strconv.FormatInt(val, 10)
		w.Write([]byte(valStr))
	}

	w.WriteHeader(http.StatusOK)
}
