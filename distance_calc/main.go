package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mcgtrt/toll-tracker/aggregator/client"
)

const kafkaTopic = "obudata"

func main() {
	var (
		serv           = NewCalcService()
		grpcListenAddr = os.Getenv("AGG_GRPC_ENDPOINT")
	)
	serv = NewLogMiddleware(serv)
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

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
