package main

import (
	"time"

	"github.com/mcgtrt/toll-tracker/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) CalculateInvoice(id int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		var (
			distance float64
			amount   float64
		)
		if inv != nil {
			distance = inv.TotalDistance
			amount = inv.Amount
		}
		logrus.WithFields(logrus.Fields{
			"took":      time.Since(start),
			"totalDist": distance,
			"obuid":     id,
			"amount":    amount,
			"func":      "CalculateInvoice",
		}).Info()
	}(time.Now())
	inv, err = m.next.CalculateInvoice(id)
	return inv, err
}

func (m *LogMiddleware) AggregateDistance(dist types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"id":   dist.OBUID,
			"dist": dist,
			"err":  err,
		}).Info("AggregateDistance")
	}(time.Now())
	err = m.next.AggregateDistance(dist)
	return err
}
