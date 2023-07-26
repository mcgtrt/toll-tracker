package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/mcgtrt/toll-tracker/types"
)

func main() {
	dr, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("[RECEIVER] Starting server")
	http.HandleFunc("/ws", dr.handleWS)
	log.Fatal(http.ListenAndServe(os.Getenv("RECEIVER_PRODUCER_ENDPOINT"), nil))
}

type DataReceiver struct {
	conn *websocket.Conn
	prod DataProducer
}

func (dr *DataReceiver) produceData(data types.OBUData) error {
	return dr.prod.ProduceData(data)
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn

	go dr.wsLoop()
}

func (dr *DataReceiver) wsLoop() {
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Fatal(err)
			continue
		}
		if err := dr.produceData(data); err != nil {
			log.Fatal(err)
		}
	}
}

func NewDataReceiver() (*DataReceiver, error) {
	var (
		p          DataProducer
		err        error
		kafkaTopic = "obudata"
	)

	p, err = NewKafkaProducer(kafkaTopic)
	if err != nil {
		return nil, err
	}

	p = NewLogMiddleware(p)

	return &DataReceiver{
		prod: p,
	}, nil
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
