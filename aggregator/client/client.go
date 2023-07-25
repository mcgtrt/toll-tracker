package client

import (
	"context"

	"github.com/mcgtrt/toll-tracker/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregateRequest) error
	GetInvoice(context.Context, int) (*types.Invoice, error)
}
