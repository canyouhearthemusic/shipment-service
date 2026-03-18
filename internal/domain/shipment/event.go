package shipment

import "time"

type Event struct {
	ID         string
	ShipmentID string
	Status     Status
	OccurredAt time.Time
}
