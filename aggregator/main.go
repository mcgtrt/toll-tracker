package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/mcgtrt/toll-tracker/types"
)

func main() {
	var (
		listenAddr = flag.String("listenAddr", ":3000", "listen address for aggregator HTTP server")
		store      = NewMemoryStore()
		service    = NewInvoiceAggregator(store)
	)
	flag.Parse()
	service = NewLogMiddleware(service)
	makeHTTPTransport(*listenAddr, service)
}

func makeHTTPTransport(listenAddr string, service Aggregator) {
	fmt.Printf("HTTP Transport running on port %s\n", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(service))
	http.HandleFunc("/invoice", handleInvoice(service))
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func handleInvoice(service Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["obu"]
		if !ok {
			writeJSON(w, http.StatusBadRequest, errJSON(fmt.Errorf("obu not provided")))
			return
		}
		obuid, err := strconv.Atoi(values[0])
		if err != nil {
			writeJSON(w, http.StatusBadRequest, errJSON(err))
			return
		}
		inv, err := service.CalculateInvoice(obuid)
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
