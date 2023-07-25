package client

import (
	"github.com/mcgtrt/toll-tracker/types"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	Endpoint string
	types.AggregatorClient
}

func NewGRPCClient(endpoint string) (*GRPCClient, error) {
	conn, err := grpc.Dial(endpoint, []grpc.DialOption{grpc.WithInsecure()}...)
	if err != nil {
		return nil, err
	}
	c := types.NewAggregatorClient(conn)
	return &GRPCClient{
		Endpoint:         endpoint,
		AggregatorClient: c,
	}, nil
}
