package aggservice

import (
	"context"

	"github.com/mcgtrt/toll-tracker/types"
)

type Middleware func(Service) Service

func New() Service {

	var service Service
	{
		service = newBasicService(NewMemoryStore())
		service = newLoggingMiddleware()(service)
		service = newinstrumentationMiddleware()(service)
	}

	return service
}

type loggingMiddleware struct {
	next Service
}

func newLoggingMiddleware() Middleware {
	return func(next Service) Service {
		return loggingMiddleware{
			next: next,
		}
	}
}

func (m loggingMiddleware) Aggregate(ctx context.Context, dist types.Distance) error {
	return nil
}

func (m loggingMiddleware) Calculate(ctx context.Context, id int) (*types.Invoice, error) {
	return nil, nil
}

type instrumentationMiddleware struct {
	next Service
}

func newinstrumentationMiddleware() Middleware {
	return func(next Service) Service {
		return instrumentationMiddleware{
			next: next,
		}
	}
}

func (m instrumentationMiddleware) Aggregate(ctx context.Context, dist types.Distance) error {
	return nil
}

func (m instrumentationMiddleware) Calculate(ctx context.Context, id int) (*types.Invoice, error) {
	return nil, nil
}
