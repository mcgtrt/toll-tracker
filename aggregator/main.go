package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

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
	log.Fatal(http.ListenAndServe(listenAddr, nil))
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
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func errJSON(err error) map[string]string {
	return map[string]string{"err": err.Error()}
}
