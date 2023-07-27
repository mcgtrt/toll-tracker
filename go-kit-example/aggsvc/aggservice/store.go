package aggservice

import (
	"fmt"

	"github.com/mcgtrt/toll-tracker/types"
)

type Storer interface {
	Get(int) (float64, error)
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

func (s *MemoryStore) Get(id int) (float64, error) {
	dist, ok := s.data[id]
	if !ok {
		return 0, fmt.Errorf("invalid obuid [%d]", id)
	}
	return dist, nil
}

func (s *MemoryStore) Insert(dist types.Distance) error {
	s.data[dist.OBUID] += dist.Value
	return nil
}
