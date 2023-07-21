package main

import (
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/mcgtrt/toll-tracker/aggregator/client"
	"github.com/mcgtrt/toll-tracker/types"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService CalcServicer
	aggrClient  *client.Client
}

func NewKafkaConsumer(topic string, calcS CalcServicer, client *client.Client) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{topic}, nil)
	return &KafkaConsumer{
		consumer:    c,
		calcService: calcS,
		aggrClient:  client,
	}, nil
}

func (c *KafkaConsumer) Start() {
	logrus.Info("kafka transport started")
	c.isRunning = true
	c.readMessageLoop()
}

func (c *KafkaConsumer) readMessageLoop() {
	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("kafka consumer error: %s", err)
			continue
		}
		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("JSON serialisation error: %s", err)
			continue
		}
		calculatedDist := c.calcService.CalculateDistance(data)
		dist := types.Distance{
			OBUID: data.OBUID,
			Value: calculatedDist,
			Unix:  time.Now().UnixNano(),
		}
		if err := c.aggrClient.AggregateInvoice(dist); err != nil {
			logrus.Error("aggregation error:", err)
			continue
		}
	}
}
