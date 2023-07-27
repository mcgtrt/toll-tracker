package aggservice

import (
	"context"
	"time"

	"github.com/go-kit/log"
	"github.com/mcgtrt/toll-tracker/types"
)

type Middleware func(Service) Service

type loggingMiddleware struct {
	next Service
	log  log.Logger
}

func newLoggingMiddleware(log log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{
			next: next,
			log:  log,
		}
	}
}

func (m loggingMiddleware) Aggregate(ctx context.Context, dist types.Distance) (err error) {
	defer func(start time.Time) {
		m.log.Log(
			"middleware", "log/aggregate",
			"took", time.Since(start),
			"dist", dist.Value,
			"obuid", dist.OBUID,
			"err", err,
		)
	}(time.Now())

	err = m.next.Aggregate(ctx, dist)
	return err
}

func (m loggingMiddleware) Calculate(ctx context.Context, id int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		m.log.Log(
			"middleware", "log/calculate",
			"took", time.Since(start),
			"obuid", id,
			"amount", inv.Amount,
			"total dist", inv.TotalDistance,
			"err", err,
		)
	}(time.Now())

	inv, err = m.next.Calculate(ctx, id)
	return inv, err
}

type instrumentationMiddleware struct {
	next Service
	log  log.Logger
}

func newinstrumentationMiddleware(log log.Logger) Middleware {
	return func(next Service) Service {
		return instrumentationMiddleware{
			next: next,
			log:  log,
		}
	}
}

func (m instrumentationMiddleware) Aggregate(ctx context.Context, dist types.Distance) (err error) {
	defer func(start time.Time) {
		m.log.Log(
			"middleware", "instrumentation/aggregate",
			"took", time.Since(start),
			"dist", dist.Value,
			"obuid", dist.OBUID,
			"err", err,
		)
	}(time.Now())

	err = m.next.Aggregate(ctx, dist)
	return err
}

func (m instrumentationMiddleware) Calculate(ctx context.Context, id int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		m.log.Log(
			"middleware", "instrumentation/calculate",
			"took", time.Since(start),
			"obuid", id,
			"amount", inv.Amount,
			"total dist", inv.TotalDistance,
			"err", err,
		)
	}(time.Now())

	inv, err = m.next.Calculate(ctx, id)
	return inv, err
}
