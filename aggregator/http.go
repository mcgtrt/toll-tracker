package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mcgtrt/toll-tracker/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type HTTPFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func (e APIError) Error() string {
	return e.Err
}

type HTTPMetricHandler struct {
	reqCounter prometheus.Counter
	reqLatency prometheus.Histogram
	errCounter prometheus.Counter
}

func newHTTPMetricHandler(reqName string) *HTTPMetricHandler {
	reqCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_request_counter", reqName),
		Name:      "aggregator",
	})
	reqLatency := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: fmt.Sprintf("http_%s_request_latency", reqName),
		Name:      "aggregator",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	errCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_error_counter", reqName),
		Name:      "aggregator",
	})
	return &HTTPMetricHandler{
		reqLatency: reqLatency,
		reqCounter: reqCounter,
		errCounter: errCounter,
	}
}

func (h *HTTPMetricHandler) instrument(next HTTPFunc) HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var err error

		defer func(start time.Time) {
			latency := time.Since(start).Seconds()
			h.reqLatency.Observe(latency)
			h.reqCounter.Inc()
			logrus.WithFields(logrus.Fields{
				"latency": latency,
				"request": r.RequestURI,
				"error":   err,
			}).Info()

			if err != nil {
				h.errCounter.Inc()
			}
		}(time.Now())

		err = next(w, r)
		return err
	}
}

func makeHTTPHandlerFunc(fn HTTPFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			apiErr, ok := err.(APIError)
			if ok {
				writeJSON(w, apiErr.Code, apiErr)
				return
			}
			writeJSON(w, http.StatusInternalServerError, err)
		}
	}
}

func handleInvoice(service Aggregator) HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != "GET" {
			return APIError{
				Code: http.StatusMethodNotAllowed,
				Err:  "method not allowed",
			}
		}
		values, ok := r.URL.Query()["obu"]
		if !ok {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  "missing obu id",
			}
		}
		obuid, err := strconv.Atoi(values[0])
		if err != nil {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  "invalid id",
			}
		}
		inv, err := service.CalculateInvoice(obuid)
		if err != nil {
			return APIError{
				Code: http.StatusInternalServerError,
				Err:  err.Error(),
			}
		}
		return writeJSON(w, http.StatusOK, inv)
	}
}

func handleAggregate(service Aggregator) HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var dist types.Distance
		if err := json.NewDecoder(r.Body).Decode(&dist); err != nil {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  "bad request",
			}
		}
		if err := service.AggregateDistance(dist); err != nil {
			return APIError{
				Code: http.StatusInternalServerError,
				Err:  err.Error(),
			}
		}
		return writeJSON(w, http.StatusOK, map[string]string{"aggregate": "ok"})
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
