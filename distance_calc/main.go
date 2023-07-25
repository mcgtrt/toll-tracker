package main

import (
	"log"

	"github.com/mcgtrt/toll-tracker/aggregator/client"
)

const (
	kafkaTopic     = "obudata"
	httpListenAddr = ":3000"
	grpcListenAddr = ":3001"
)

func main() {
	serv := NewCalcService()
	serv = NewLogMiddleware(serv)
	// httpClient := client.NewHTTPClient(httpListenAddr)
	grpcClient, err := client.NewGRPCClient(grpcListenAddr)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, serv, grpcClient)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
