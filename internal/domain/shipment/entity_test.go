package shipment_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/canyouhearthemusic/shipment-service/internal/domain/shipment"
)

func newTestShipment(t *testing.T) *shipment.Shipment {
	t.Helper()
	s, err := shipment.NewShipment(
		"REF-001",
		"NYC",
		"LA",
		"John",
		"T1",
		100,
		70,
	)
	require.NoError(t, err)

	return s
}

func TestNewShipment(t *testing.T) {
	t.Run("valid inputs", func(t *testing.T) {
		s, err := shipment.NewShipment(
			"REF-001",
			"NYC",
			"LA",
			"John",
			"T1",
			100,
			70,
		)
		require.NoError(t, err)

		assert.NotEmpty(t, s.ID)
		assert.Equal(t, "REF-001", s.ReferenceNumber)
		assert.Equal(t, "NYC", s.Origin)
		assert.Equal(t, "LA", s.Destination)
		assert.Equal(t, "John", s.DriverName)
		assert.Equal(t, "T1", s.UnitNumber)
		assert.Equal(t, float64(100), s.Amount)
		assert.Equal(t, float64(70), s.DriverRevenue)
		assert.Equal(t, shipment.StatusPending, s.CurrentStatus)
		assert.False(t, s.CreatedAt.IsZero())
		assert.False(t, s.UpdatedAt.IsZero())
	})

	t.Run("initial event recorded", func(t *testing.T) {
		s := newTestShipment(t)

		require.Len(t, s.Events, 1)
		assert.NotEmpty(t, s.Events[0].ID)
		assert.Equal(t, s.ID, s.Events[0].ShipmentID)
		assert.Equal(t, shipment.StatusPending, s.Events[0].Status)
	})

	t.Run("missing required fields", func(t *testing.T) {
		tests := []struct {
			name              string
			ref, origin, dest string
		}{
			{"missing reference", "", "A", "B"},
			{"missing origin", "REF", "", "B"},
			{"missing destination", "REF", "A", ""},
			{"all empty", "", "", ""},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := shipment.NewShipment(tt.ref, tt.origin, tt.dest, "D", "U", 1, 1)
				assert.ErrorIs(t, err, shipment.ErrInvalidInput)
			})
		}
	})
}
