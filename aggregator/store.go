package main

import "github.com/mcgtrt/toll-tracker/types"

type Storer interface {
	Insert(types.Distance) error
}

type MemoryStore struct {
	data map[int]float64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[int]float64),
	}
}

func (s *MemoryStore) Insert(dist types.Distance) error {
	s.data[dist.OBUID] += dist.Value
	return nil
}
