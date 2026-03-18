package shipment

import (
	"context"
	"fmt"

	domain "github.com/canyouhearthemusic/shipment-service/internal/domain/shipment"
)

type Service interface {
	CreateShipment(ctx context.Context, req CreateShipmentRequest) (*domain.Shipment, error)
	GetShipment(ctx context.Context, id string) (*domain.Shipment, error)
	AddShipmentEvent(ctx context.Context, id string, status domain.Status) (*domain.Shipment, error)
	GetShipmentEvents(ctx context.Context, id string) ([]domain.Event, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateShipment(ctx context.Context, req CreateShipmentRequest) (*domain.Shipment, error) {
	shipment, err := domain.NewShipment(
		req.ReferenceNumber,
		req.Origin,
		req.Destination,
		req.DriverName,
		req.UnitNumber,
		req.Amount,
		req.DriverRevenue,
	)
	if err != nil {
		return nil, fmt.Errorf("create shipment: %w", err)
	}

	if err := s.repo.Save(ctx, shipment); err != nil {
		return nil, fmt.Errorf("save shipment: %w", err)
	}

	return shipment, nil
}

func (s *service) GetShipment(ctx context.Context, id string) (*domain.Shipment, error) {
	shipment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get shipment %s: %w", id, err)
	}

	return shipment, nil
}

func (s *service) AddShipmentEvent(ctx context.Context, id string, status domain.Status) (*domain.Shipment, error) {
	shipment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find shipment %s: %w", id, err)
	}

	if err := shipment.AddEvent(status); err != nil {
		return nil, fmt.Errorf("add event: %w", err)
	}

	if err := s.repo.Save(ctx, shipment); err != nil {
		return nil, fmt.Errorf("save shipment: %w", err)
	}

	return shipment, nil
}

func (s *service) GetShipmentEvents(ctx context.Context, id string) ([]domain.Event, error) {
	shipment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find shipment %s: %w", id, err)
	}

	return shipment.Events, nil
}
