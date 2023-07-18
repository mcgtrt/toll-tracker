package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mcgtrt/toll-tracker/types"
)

func main() {
	dr := DataReceiver{
		msg: make(chan types.OBUData, 128),
	}
	http.HandleFunc("/ws", dr.handleWS)
	http.ListenAndServe(":3000", nil)
}

type DataReceiver struct {
	msg  chan types.OBUData
	conn *websocket.Conn
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
	fmt.Println("[RECEIVER] Successfully connected with client!")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Fatal(err)
			continue
		}
		fmt.Printf("ID [%d] :: <lat:'%.2f', long:'%.2f'>\n", data.OBUID, data.Lat, data.Long)
		dr.msg <- data
	}
}
