package zoo

import (
	"context"
	"errors"
	"sync"
)

type MemoryStore struct {
	mu       sync.RWMutex
	routes   []RoutePlan
	failWith error
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (s *MemoryStore) Save(ctx context.Context, plan RoutePlan) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.failWith != nil {
		return s.failWith
	}
	s.routes = append(s.routes, plan)
	return nil
}

func (s *MemoryStore) SetFailure(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.failWith = err
}

func (s *MemoryStore) Routes() []RoutePlan {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]RoutePlan, len(s.routes))
	copy(result, s.routes)
	return result
}

func (s *MemoryStore) ClearFailure() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.failWith = nil
}

var ErrStoreUnavailable = errors.New("route store unavailable")
