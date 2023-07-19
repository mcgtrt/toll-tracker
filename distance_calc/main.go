package main

import "log"

const kafkaTopic = "obudata"

func main() {
	serv := NewCalcService()
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, serv)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
