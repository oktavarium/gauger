package handlers

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

type Handler struct {
	storage metricsSaver
}

func NewHandler(storage metricsSaver) *Handler {
	return &Handler{
		storage: storage,
	}
}

func (h *Handler) GetHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(h.storage.GetAll()))
}

func (h *Handler) UpdateHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
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
		h.storage.SaveGauge(metricName, val)

	} else {
		val, err := strconv.ParseInt(metricValueStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.storage.UpdateCounter(metricName, val)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ValueHandle(w http.ResponseWriter, r *http.Request) {
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
		val, ok := h.storage.GetGauger(metricName)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		valStr := strconv.FormatFloat(val, 'f', -1, 64)
		w.Write([]byte(valStr))
	} else {
		val, ok := h.storage.GetCounter(metricName)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		valStr := strconv.FormatInt(val, 10)
		w.Write([]byte(valStr))
	}
}
