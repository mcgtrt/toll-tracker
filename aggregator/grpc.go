package main

import (
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

func (s *GRPCServer) AggregateDistance(req *types.AggregateRequest) error {
	dist := types.Distance{
		OBUID: int(req.OBUID),
		Value: req.Value,
		Unix:  req.Unix,
	}
	return s.service.AggregateDistance(dist)
}
