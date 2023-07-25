package main

import (
	"fmt"

	"github.com/mcgtrt/toll-tracker/types"
)

// As this application is a demo of how to structure, connect, and build transport
// between microservices, it does not include the business logic or integration
// with the database that should be handled in this file (e.g. MongoDB would be good
// for that use case). All operations are handled in memory to prove the concept and
// make the application work properly.
type Storer interface {
	GetDistanceByOBUID(int) (float64, error)
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

func (s *MemoryStore) GetDistanceByOBUID(id int) (float64, error) {
	dist, ok := s.data[id]
	if !ok {
		return 0, fmt.Errorf("invalid obuid [%d]", id)
	}
	return dist, nil
}

func (s *MemoryStore) Insert(dist types.Distance) error {
	s.data[dist.OBUID] += dist.Value
	fmt.Printf("Inserting ID: %d\nNew distance: %.2f\n", dist.OBUID, s.data[dist.OBUID])
	return nil
}
