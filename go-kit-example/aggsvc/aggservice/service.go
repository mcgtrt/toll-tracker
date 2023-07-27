package aggservice

import (
	"context"

	"github.com/mcgtrt/toll-tracker/types"
)

const basePrice = 0.12

type Service interface {
	Aggregate(context.Context, types.Distance) error
	Calculate(context.Context, int) (*types.Invoice, error)
}

func New() Service {

	var service Service
	{
		service = newBasicService(NewMemoryStore())
		service = newLoggingMiddleware()(service)
		service = newinstrumentationMiddleware()(service)
	}

	return service
}

type BasicService struct {
	store Storer
}

func newBasicService(store Storer) Service {
	return &BasicService{
		store: store,
	}
}

func (s *BasicService) Aggregate(_ context.Context, dist types.Distance) error {
	return s.store.Insert(dist)
}

func (s *BasicService) Calculate(_ context.Context, obuid int) (*types.Invoice, error) {
	dist, err := s.store.Get(obuid)
	if err != nil {
		return nil, err
	}
	inv := &types.Invoice{
		OBUID:         obuid,
		TotalDistance: dist,
		Amount:        basePrice * dist,
	}
	return inv, nil
}
