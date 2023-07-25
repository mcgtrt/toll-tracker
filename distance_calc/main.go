package main

import (
	"log"

	"github.com/mcgtrt/toll-tracker/aggregator/client"
)

const (
	kafkaTopic          = "obudata"
	aggregationEndpoint = "http://127.0.0.1:3000/aggregate"
)

func main() {
	serv := NewCalcService()
	serv = NewLogMiddleware(serv)
	grpcClient, err := client.NewGRPCClient(aggregationEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, serv, grpcClient)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
