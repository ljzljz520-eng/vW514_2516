package zoo

import (
	"context"
	"fmt"
)

type Service struct {
	planner *Planner
	store   RouteStore
}

func NewService(planner *Planner, store RouteStore) *Service {
	return &Service{planner: planner, store: store}
}

func (s *Service) Execute(ctx context.Context, req RouteRequest) (RoutePlan, error) {
	if s == nil || s.planner == nil || s.store == nil {
		return RoutePlan{}, fmt.Errorf("%w: planner and store are required", ErrInvalidRequest)
	}
	if err := s.planner.Validate(req); err != nil {
		return RoutePlan{}, err
	}
	plan, err := s.planner.Plan(req)
	if err != nil {
		return RoutePlan{}, err
	}
	if err := s.store.Save(ctx, plan); err != nil {
		return RoutePlan{}, err
	}
	return plan, nil
}

type AdminService struct {
	catalog *Catalog
}

func NewAdminService(catalog *Catalog) *AdminService {
	return &AdminService{catalog: catalog}
}

func (a *AdminService) UpsertNode(node Node) error {
	if a == nil || a.catalog == nil {
		return fmt.Errorf("%w: catalog is required", ErrInvalidRequest)
	}
	return a.catalog.UpsertNode(node)
}

func (a *AdminService) AddConnection(connection Connection) error {
	if a == nil || a.catalog == nil {
		return fmt.Errorf("%w: catalog is required", ErrInvalidRequest)
	}
	return a.catalog.AddConnection(connection)
}

func (a *AdminService) SetCrowd(from, to string, crowd int) error {
	if a == nil || a.catalog == nil {
		return fmt.Errorf("%w: catalog is required", ErrInvalidRequest)
	}
	return a.catalog.SetCrowd(from, to, crowd)
}
