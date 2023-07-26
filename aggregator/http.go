package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mcgtrt/toll-tracker/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type HTTPMetricHandler struct {
	reqCounter prometheus.Counter
}

func newHTTPMetricHandler(reqName string) *HTTPMetricHandler {
	reqCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_request_counter", reqName),
		Name:      "aggregator",
	})
	return &HTTPMetricHandler{
		reqCounter: reqCounter,
	}
}

func (h *HTTPMetricHandler) instrument(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.reqCounter.Inc()
		next(w, r)
	}
}

func handleInvoice(service Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["obu"]
		if !ok {
			writeJSON(w, http.StatusBadRequest, errJSON(fmt.Errorf("obu not provided")))
			return
		}
		obuid, err := strconv.Atoi(values[0])
		fmt.Printf("OBUID: %d\n", obuid)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, errJSON(err))
			return
		}
		inv, err := service.CalculateInvoice(obuid)
		fmt.Printf("INV: %v, ERR: %v", inv, err)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, errJSON(err))
			return
		}
		writeJSON(w, http.StatusOK, inv)
	}
}

func handleAggregate(service Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dist types.Distance
		if err := json.NewDecoder(r.Body).Decode(&dist); err != nil {
			writeJSON(w, http.StatusBadRequest, errJSON(err))
			return
		}
		if err := service.AggregateDistance(dist); err != nil {
			writeJSON(w, http.StatusInternalServerError, errJSON(err))
			return
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func errJSON(err error) map[string]string {
	return map[string]string{"err": err.Error()}
}
