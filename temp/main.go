package main

import (
	"context"
	"log"

	"github.com/mcgtrt/toll-tracker/aggregator/client"
	"github.com/mcgtrt/toll-tracker/types"
)

func main() {
	c, err := client.NewGRPCClient(":3001")
	if err != nil {
		log.Fatal(err)
	}
	if err := c.Aggregate(context.Background(), &types.AggregateRequest{
		OBUID: 1,
		Value: 2,
		Unix:  3,
	}); err != nil {
		log.Fatal(err)
	}
}
