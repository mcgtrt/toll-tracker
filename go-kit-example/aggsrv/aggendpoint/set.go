package aggendpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/mcgtrt/toll-tracker/go-kit-example/aggsrv/aggservice"
	"github.com/mcgtrt/toll-tracker/types"
)

type Set struct {
	AggregateEndpoint endpoint.Endpoint
	CalculateEndpoint endpoint.Endpoint
}

func (s Set) Aggregate(ctx context.Context, dist types.Distance) error {
	res, err := s.AggregateEndpoint(ctx, AggregateRequest{
		OBUID: dist.OBUID,
		Value: dist.Value,
		Unix:  dist.Unix,
	})
	if err != nil {
		return err
	}
	ares := res.(AggregateResponse)
	return ares.Err
}

func (s Set) Calculate(ctx context.Context, id int) (*types.Invoice, error) {
	res, err := s.CalculateEndpoint(ctx, CalculateRequest{
		OBUID: id,
	})
	if err != nil {
		return nil, err
	}
	cresp := res.(CalculateResponse)
	return &types.Invoice{
		OBUID:         cresp.OBUID,
		TotalDistance: cresp.TotalDistance,
		Amount:        cresp.Amount,
	}, cresp.Err
}

func MakeAggretageEndpoint(s aggservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AggregateRequest)
		err = s.Aggregate(ctx, types.Distance{
			OBUID: req.OBUID,
			Value: req.Value,
			Unix:  req.Unix,
		})
		return AggregateResponse{Err: err}, nil
	}
}

func MakeCalculateEndpoint(s aggservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CalculateRequest)
		inv, err := s.Calculate(ctx, req.OBUID)
		return CalculateResponse{
			OBUID:         inv.OBUID,
			TotalDistance: inv.TotalDistance,
			Amount:        inv.Amount,
			Err:           err,
		}, nil
	}
}

type AggregateRequest struct {
	OBUID int     `json:"obuID"`
	Value float64 `json:"value"`
	Unix  int64   `jsin:"unix"`
}

type AggregateResponse struct {
	Err error `json:"err"`
}

type CalculateRequest struct {
	OBUID int `json:"obuID"`
}

type CalculateResponse struct {
	OBUID         int     `json:"obuID"`
	TotalDistance float64 `json:"totalDistance"`
	Amount        float64 `json:"amount"`
	Err           error   `json:"err"`
}
