package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mcgtrt/toll-tracker/types"
)

func main() {
	dr, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}
	defer dr.prod.CloseProducer()

	fmt.Println("[RECEIVER] Starting the server")
	http.HandleFunc("/ws", dr.handleWS)
	http.ListenAndServe(":3000", nil)
}

type DataReceiver struct {
	conn *websocket.Conn
	prod DataProducer
}

func NewDataReceiver() (*DataReceiver, error) {
	prod, err := NewKafkaProducer()
	if err != nil {
		return nil, err
	}
	return &DataReceiver{
		prod: prod,
	}, nil
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

	go dr.wsReceiveLoop()
}

func (dr *DataReceiver) wsReceiveLoop() {
	fmt.Println("[RECEIVER] Started receiving channel")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Fatal(err)
			continue
		}
		if err := dr.prod.ProduceData(data); err != nil {
			fmt.Println("kafka produced an error:", err)
		}
		fmt.Println("[RECEIVER] Received ", data)
	}
}
