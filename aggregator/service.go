package main

import (
	"github.com/mcgtrt/toll-tracker/types"
)

type Aggregator interface {
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

func (a InvoiceAggregator) AggregateDistance(dist types.Distance) error {
	return a.store.Insert(dist)
}
