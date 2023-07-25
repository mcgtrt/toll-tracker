package main

import (
	"context"

	"github.com/mcgtrt/toll-tracker/types"
)

type GRPCServer struct {
	types.UnimplementedAggregatorServer
	service Aggregator
}

func NewGRPCServer(service Aggregator) *GRPCServer {
	return &GRPCServer{
		service: service,
	}
}

// This protoc method returns *types.None as in this example we're not listening for any return,
// only if aggregation was successful
func (s *GRPCServer) Aggregate(ctx context.Context, req *types.AggregateRequest) (*types.None, error) {
	dist := types.Distance{
		OBUID: int(req.OBUID),
		Value: req.Value,
		Unix:  req.Unix,
	}
	return &types.None{}, s.service.AggregateDistance(dist)
}
