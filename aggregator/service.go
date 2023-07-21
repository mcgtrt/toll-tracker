package main

import (
	"fmt"

	"github.com/mcgtrt/toll-tracker/types"
)

const basePrice = 0.12

type Aggregator interface {
	CalculateInvoice(int) (*types.Invoice, error)
	AggregateDistance(types.Distance) error
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
	return InvoiceAggregator{
		store: store,
	}
}

func (a InvoiceAggregator) CalculateInvoice(id int) (*types.Invoice, error) {
	dist, err := a.store.GetDistanceByOBUID(id)
	if err != nil {
		return nil, err
	}
	inv := &types.Invoice{
		OBUID:         id,
		TotalDistance: dist,
		Amount:        basePrice * dist,
	}
	fmt.Println(inv)
	return inv, nil
}

func (a InvoiceAggregator) AggregateDistance(dist types.Distance) error {
	return a.store.Insert(dist)
}
