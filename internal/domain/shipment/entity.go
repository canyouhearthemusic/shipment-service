package shipment

import (
	"time"

	"github.com/google/uuid"
)

type Shipment struct {
	ID              string
	ReferenceNumber string
	Origin          string
	Destination     string
	CurrentStatus   Status
	DriverName      string
	UnitNumber      string
	Amount          float64
	DriverRevenue   float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Events          []Event
}

func NewShipment(
	referenceNumber, origin, destination, driverName, unitNumber string,
	amount, driverRevenue float64,
) (*Shipment, error) {
	if referenceNumber == "" || origin == "" || destination == "" {
		return nil, ErrInvalidInput
	}

	now := time.Now().UTC()
	s := &Shipment{
		ID:              uuid.New().String(),
		ReferenceNumber: referenceNumber,
		Origin:          origin,
		Destination:     destination,
		CurrentStatus:   StatusPending,
		DriverName:      driverName,
		UnitNumber:      unitNumber,
		Amount:          amount,
		DriverRevenue:   driverRevenue,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	s.Events = []Event{{
		ID:         uuid.New().String(),
		ShipmentID: s.ID,
		Status:     StatusPending,
		OccurredAt: now,
	}}

	return s, nil
}

func (s *Shipment) AddEvent(status Status) error {
	if !s.CurrentStatus.CanTransitionTo(status) {
		return ErrInvalidTransition
	}

	now := time.Now().UTC()
	s.Events = append(s.Events, Event{
		ID:         uuid.New().String(),
		ShipmentID: s.ID,
		Status:     status,
		OccurredAt: now,
	})
	s.CurrentStatus = status
	s.UpdatedAt = now

	return nil
}
