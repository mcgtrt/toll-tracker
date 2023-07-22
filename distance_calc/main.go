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
	var (
		client = client.NewHTTPClient(aggregationEndpoint)
		serv   = NewCalcService()
	)
	serv = NewLogMiddleware(serv)
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, serv, client)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
