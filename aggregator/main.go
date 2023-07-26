package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/mcgtrt/toll-tracker/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
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
	var (
		aggMetricHandler = newHTTPMetricHandler("aggregate")
		invMetricHandler = newHTTPMetricHandler("invoice")
	)

	logrus.Infof("HTTP Transport running on port %s\n", listenAddr)

	http.HandleFunc("/aggregate", aggMetricHandler.instrument(handleAggregate(service)))
	http.HandleFunc("/invoice", invMetricHandler.instrument(handleInvoice(service)))
	http.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(listenAddr, nil)
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
