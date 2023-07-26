package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/mcgtrt/toll-tracker/types"
	"github.com/sirupsen/logrus"
)

var (
	generateWaitTime = time.Second * 5
	wsEndpoint       = fmt.Sprintf("ws://127.0.0.1%s/ws", os.Getenv("RECEIVER_PRODUCER_ENDPOINT"))
)

// This service simulates sending real world OBUs(On Board Units) data
// that will be later received by another service for processing
func main() {
	obuIDs := generateOBUIDs(20)
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		for i := 0; i < len(obuIDs); i++ {
			lat, long := generateLatLong()
			data := types.OBUData{
				OBUID: obuIDs[i],
				Lat:   lat,
				Long:  long,
			}
			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
			logrus.WithFields(logrus.Fields{
				"id":   data.OBUID,
				"lat":  data.Lat,
				"long": data.Long,
			}).Info("Generating Random OBU Data")
		}
		time.Sleep(generateWaitTime)
	}
}

func generateLatLong() (float64, float64) {
	return generateCoord(), generateCoord()
}

func generateCoord() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()
	return n + f
}

func generateOBUIDs(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(math.MaxInt)
	}
	return ids
}

func init() {
	rand.Seed(time.Now().UnixNano())
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
