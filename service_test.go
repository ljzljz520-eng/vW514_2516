package zoo

import (
	"context"
	"errors"
	"testing"
)

type failingStore struct {
	err   error
	calls int
}

func (s *failingStore) Save(context.Context, RoutePlan) error {
	s.calls++
	return s.err
}

func TestExecuteReturnsPersistenceError(t *testing.T) {
	sentinel := errors.New("fixture persistence failure")
	store := &failingStore{err: sentinel}
	service := NewService(NewPlanner(FixedFixture()), store)
	_, err := service.Execute(context.Background(), RouteRequest{EntranceID: "entrance-north", AnimalAreaIDs: []string{"tiger"}, LunchID: "picnic-lawn", ExitID: "exit-east", StartMinute: 540, Stroller: true})
	if !errors.Is(err, sentinel) {
		t.Fatalf("got error %v, want %v", err, sentinel)
	}
	if store.calls != 1 {
		t.Fatalf("save calls = %d", store.calls)
	}
}

func TestExecuteRejectsInvalidRequestBeforeSave(t *testing.T) {
	store := &failingStore{err: ErrStoreUnavailable}
	service := NewService(NewPlanner(FixedFixture()), store)
	_, err := service.Execute(context.Background(), RouteRequest{EntranceID: "missing", AnimalAreaIDs: []string{"tiger"}, LunchID: "picnic-lawn", ExitID: "exit-east", StartMinute: 540})
	if !errors.Is(err, ErrUnknownNode) {
		t.Fatalf("got error %v", err)
	}
	if store.calls != 0 {
		t.Fatalf("save calls = %d", store.calls)
	}
}

func TestMemoryStorePersistsSuccessfulPlan(t *testing.T) {
	store := NewMemoryStore()
	service := NewService(NewPlanner(FixedFixture()), store)
	_, err := service.Execute(context.Background(), RouteRequest{EntranceID: "entrance-south", AnimalAreaIDs: []string{"savanna"}, LunchID: "family-cafe", ExitID: "exit-west", StartMinute: 600, Stroller: false})
	if err != nil {
		t.Fatalf("execute failed: %v", err)
	}
	if len(store.Routes()) != 1 {
		t.Fatalf("routes = %d", len(store.Routes()))
	}
}
