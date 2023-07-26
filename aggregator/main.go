package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/mcgtrt/toll-tracker/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	var (
		httpListenAddr = os.Getenv("AGG_HTTP_ENDPOINT")
		grpcListenAddr = os.Getenv("AGG_GRPC_ENDPOINT")
		store          = NewMemoryStore()
		service        = NewInvoiceAggregator(store)
	)

	service = NewMetricsMiddleware(service)
	service = NewLogMiddleware(service)

	go func() {
		log.Fatal(makeGRPCTransport(grpcListenAddr, service))
	}()
	log.Fatal(makeHTTPTransport(httpListenAddr, service))
}

func makeGRPCTransport(listenAddr string, srv Aggregator) error {
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer func() {
		fmt.Println("Closing gRPC server")
		ln.Close()
	}()

	server := grpc.NewServer([]grpc.ServerOption{}...)
	types.RegisterAggregatorServer(server, NewGRPCServer(srv))
	fmt.Printf("GRPC Transport running on port %s\n", listenAddr)
	return server.Serve(ln)
}

func makeHTTPTransport(listenAddr string, service Aggregator) error {
	http.HandleFunc("/aggregate", handleAggregate(service))
	http.HandleFunc("/invoice", handleInvoice(service))
	http.Handle("/metrics", promhttp.Handler())
	fmt.Printf("HTTP Transport running on port %s\n", listenAddr)
	return http.ListenAndServe(listenAddr, nil)
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

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
