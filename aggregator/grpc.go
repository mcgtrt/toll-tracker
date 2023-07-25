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

func (s *GRPCServer) Aggregate(ctx context.Context, req *types.AggregateRequest) (*types.None, error) {
	dist := types.Distance{
		OBUID: int(req.OBUID),
		Value: req.Value,
		Unix:  req.Unix,
	}
	return &types.None{}, s.service.AggregateDistance(dist)
}
