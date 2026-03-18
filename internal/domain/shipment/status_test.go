package shipment_test

import (
	"fmt"
	"testing"

	"github.com/canyouhearthemusic/shipment-service/internal/domain/shipment"
)

// Check matrix in validTransitions map
func TestCanTransitionTo(t *testing.T) {
	tests := []struct {
		from shipment.Status
		to   shipment.Status
		want bool
	}{
		{shipment.StatusPending, shipment.StatusPickedUp, true},
		{shipment.StatusPending, shipment.StatusCancelled, true},
		{shipment.StatusPending, shipment.StatusInTransit, false},
		{shipment.StatusPending, shipment.StatusDelivered, false},

		{shipment.StatusPickedUp, shipment.StatusInTransit, true},
		{shipment.StatusPickedUp, shipment.StatusCancelled, true},
		{shipment.StatusPickedUp, shipment.StatusPending, false},
		{shipment.StatusPickedUp, shipment.StatusDelivered, false},

		{shipment.StatusInTransit, shipment.StatusDelivered, true},
		{shipment.StatusInTransit, shipment.StatusCancelled, true},
		{shipment.StatusInTransit, shipment.StatusPending, false},
		{shipment.StatusInTransit, shipment.StatusPickedUp, false},

		{shipment.StatusDelivered, shipment.StatusPending, false},
		{shipment.StatusDelivered, shipment.StatusCancelled, false},

		{shipment.StatusCancelled, shipment.StatusPending, false},
		{shipment.StatusCancelled, shipment.StatusPickedUp, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_%s", string(tt.from), string(tt.to)), func(t *testing.T) {
			got := tt.from.CanTransitionTo(tt.to)
			if got != tt.want {
				t.Errorf("(%s).CanTransitionTo(%s) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}
