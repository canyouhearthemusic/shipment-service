package memory

import (
	"context"
	"sync"

	domain "github.com/canyouhearthemusic/shipment-service/internal/domain/shipment"
)

type shipmentRepository struct {
	mu   sync.RWMutex
	data map[string]*domain.Shipment
}

func NewShipmentRepository() *shipmentRepository {
	return &shipmentRepository{
		data: make(map[string]*domain.Shipment),
	}
}

func (r *shipmentRepository) Save(_ context.Context, s *domain.Shipment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cp := *s
	events := make([]domain.Event, len(s.Events))
	copy(events, s.Events)
	cp.Events = events

	r.data[s.ID] = &cp

	return nil
}

func (r *shipmentRepository) FindByID(_ context.Context, id string) (*domain.Shipment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s, ok := r.data[id]
	if !ok {
		return nil, domain.ErrShipmentNotFound
	}

	cp := *s
	events := make([]domain.Event, len(s.Events))
	copy(events, s.Events)
	cp.Events = events

	return &cp, nil
}

func (r *shipmentRepository) FindAll(_ context.Context) ([]*domain.Shipment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]*domain.Shipment, 0, len(r.data))
	for _, s := range r.data {
		cp := *s
		events := make([]domain.Event, len(s.Events))
		copy(events, s.Events)
		cp.Events = events
		out = append(out, &cp)
	}

	return out, nil
}
