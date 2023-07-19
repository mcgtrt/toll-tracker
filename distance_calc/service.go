package main

import (
	"math"

	"github.com/mcgtrt/toll-tracker/types"
)

type CalcServicer interface {
	CalculateDistance(types.OBUData) float64
}

type CalcService struct {
	distance     map[int]float64
	previousData map[int][]float64
}

func NewCalcService() *CalcService {
	return &CalcService{
		distance:     make(map[int]float64),
		previousData: make(map[int][]float64),
	}
}

func (s CalcService) CalculateDistance(data types.OBUData) float64 {
	id := data.OBUID
	prev := s.previousData[id]
	if len(prev) == 0 {
		s.previousData[id] = []float64{data.Lat, data.Long}
		return 0
	}
	dist := math.Sqrt(math.Pow(data.Lat-prev[0], 2) + math.Pow(data.Long-prev[1], 2))
	newDist := s.distance[id] + dist
	s.distance[id] = newDist
	s.previousData[id] = []float64{data.Lat, data.Long}
	return newDist
}
