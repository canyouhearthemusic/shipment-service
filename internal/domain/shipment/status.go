package shipment

import "slices"

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusPickedUp  Status = "PICKED_UP"
	StatusInTransit Status = "IN_TRANSIT"
	StatusDelivered Status = "DELIVERED"
	StatusCancelled Status = "CANCELLED"
)

var validTransitions = map[Status][]Status{
	StatusPending:   {StatusPickedUp, StatusCancelled},
	StatusPickedUp:  {StatusInTransit, StatusCancelled},
	StatusInTransit: {StatusDelivered, StatusCancelled},
}

func (s Status) CanTransitionTo(next Status) bool {
	return slices.Contains(validTransitions[s], next)
}
