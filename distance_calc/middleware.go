package main

import (
	"time"

	"github.com/mcgtrt/toll-tracker/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next CalcServicer
}

func NewLogMiddleware(next CalcServicer) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) CalculateDistance(data types.OBUData) (dist float64) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"dist": dist,
		}).Info()
	}(time.Now())
	return m.next.CalculateDistance(data)
}
