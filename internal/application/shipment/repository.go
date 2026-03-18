package shipment

import (
	"context"

	domain "github.com/canyouhearthemusic/shipment-service/internal/domain/shipment"
)

//go:generate mockgen -destination=mocks/repository_mock.go -package=mocks . Repository
type Repository interface {
	Save(ctx context.Context, s *domain.Shipment) error
	FindByID(ctx context.Context, id string) (*domain.Shipment, error)
	FindAll(ctx context.Context) ([]*domain.Shipment, error)
}
