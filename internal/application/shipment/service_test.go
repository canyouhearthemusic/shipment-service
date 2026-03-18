package shipment_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	appshipment "github.com/canyouhearthemusic/shipment-service/internal/application/shipment"
	"github.com/canyouhearthemusic/shipment-service/internal/application/shipment/mocks"
	domain "github.com/canyouhearthemusic/shipment-service/internal/domain/shipment"
)

func TestServiceCreateShipment(t *testing.T) {
	tests := []struct {
		name  string
		mock  func(repo *mocks.MockRepository)
		req   appshipment.CreateShipmentRequest
		check func(t *testing.T, s *domain.Shipment, err error)
	}{
		{
			name: "success",
			mock: func(repo *mocks.MockRepository) {
				repo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
			},
			req: appshipment.CreateShipmentRequest{
				ReferenceNumber: "REF-001",
				Origin:          "NYC",
				Destination:     "LA",
				DriverName:      "John",
				UnitNumber:      "T1",
				Amount:          100,
				DriverRevenue:   70,
			},
			check: func(t *testing.T, s *domain.Shipment, err error) {
				require.NoError(t, err)
				assert.Equal(t, "REF-001", s.ReferenceNumber)
				assert.Equal(t, domain.StatusPending, s.CurrentStatus)
				assert.Len(t, s.Events, 1)
			},
		},
		{
			name: "invalid input",
			mock: func(repo *mocks.MockRepository) {},
			req: appshipment.CreateShipmentRequest{
				Origin:      "NYC",
				Destination: "LA",
			},
			check: func(t *testing.T, _ *domain.Shipment, err error) {
				assert.ErrorIs(t, err, domain.ErrInvalidInput)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockRepository(ctrl)
			tt.mock(repo)

			svc := appshipment.NewService(repo)
			s, err := svc.CreateShipment(context.Background(), tt.req)
			tt.check(t, s, err)
		})
	}
}

func TestServiceGetShipment(t *testing.T) {
	tests := []struct {
		name  string
		mock  func(repo *mocks.MockRepository)
		id    string
		check func(t *testing.T, s *domain.Shipment, err error)
	}{
		{
			name: "found",
			mock: func(repo *mocks.MockRepository) {
				repo.EXPECT().FindByID(gomock.Any(), "abc-123").
					Return(&domain.Shipment{ID: "abc-123", ReferenceNumber: "REF-001"}, nil)
			},
			id: "abc-123",
			check: func(t *testing.T, s *domain.Shipment, err error) {
				require.NoError(t, err)
				assert.Equal(t, "abc-123", s.ID)
			},
		},
		{
			name: "not found",
			mock: func(repo *mocks.MockRepository) {
				repo.EXPECT().FindByID(gomock.Any(), "nonexistent").
					Return(nil, domain.ErrShipmentNotFound)
			},
			id: "nonexistent",
			check: func(t *testing.T, _ *domain.Shipment, err error) {
				assert.ErrorIs(t, err, domain.ErrShipmentNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockRepository(ctrl)
			tt.mock(repo)

			svc := appshipment.NewService(repo)
			s, err := svc.GetShipment(context.Background(), tt.id)
			tt.check(t, s, err)
		})
	}
}

func TestServiceAddShipmentEvent(t *testing.T) {
	pendingShipment := &domain.Shipment{
		ID: "abc-123", CurrentStatus: domain.StatusPending,
		Events: []domain.Event{{Status: domain.StatusPending}},
	}

	tests := []struct {
		name   string
		mock   func(repo *mocks.MockRepository)
		id     string
		status domain.Status
		check  func(t *testing.T, s *domain.Shipment, err error)
	}{
		{
			name: "valid transition",
			mock: func(repo *mocks.MockRepository) {
				repo.EXPECT().FindByID(gomock.Any(), "abc-123").Return(pendingShipment, nil)
				repo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
			},
			id:     "abc-123",
			status: domain.StatusPickedUp,
			check: func(t *testing.T, s *domain.Shipment, err error) {
				require.NoError(t, err)
				assert.Equal(t, domain.StatusPickedUp, s.CurrentStatus)
			},
		},
		{
			name: "invalid transition",
			mock: func(repo *mocks.MockRepository) {
				repo.EXPECT().FindByID(gomock.Any(), "abc-123").Return(pendingShipment, nil)
			},
			id:     "abc-123",
			status: domain.StatusDelivered,
			check: func(t *testing.T, _ *domain.Shipment, err error) {
				assert.ErrorIs(t, err, domain.ErrInvalidTransition)
			},
		},
		{
			name: "not found",
			mock: func(repo *mocks.MockRepository) {
				repo.EXPECT().FindByID(gomock.Any(), "nonexistent").Return(nil, domain.ErrShipmentNotFound)
			},
			id:     "nonexistent",
			status: domain.StatusPickedUp,
			check: func(t *testing.T, _ *domain.Shipment, err error) {
				assert.ErrorIs(t, err, domain.ErrShipmentNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockRepository(ctrl)
			tt.mock(repo)

			svc := appshipment.NewService(repo)
			s, err := svc.AddShipmentEvent(context.Background(), tt.id, tt.status)
			tt.check(t, s, err)
		})
	}
}

func TestServiceGetShipmentEvents(t *testing.T) {
	events := []domain.Event{
		{Status: domain.StatusPending},
		{Status: domain.StatusPickedUp},
		{Status: domain.StatusInTransit},
	}

	tests := []struct {
		name  string
		mock  func(repo *mocks.MockRepository)
		id    string
		check func(t *testing.T, events []domain.Event, err error)
	}{
		{
			name: "returns event history",
			mock: func(repo *mocks.MockRepository) {
				repo.EXPECT().FindByID(gomock.Any(), "abc-123").
					Return(&domain.Shipment{ID: "abc-123", Events: events}, nil)
			},
			id: "abc-123",
			check: func(t *testing.T, got []domain.Event, err error) {
				require.NoError(t, err)
				require.Len(t, got, 3)
				assert.Equal(t, domain.StatusPending, got[0].Status)
				assert.Equal(t, domain.StatusPickedUp, got[1].Status)
				assert.Equal(t, domain.StatusInTransit, got[2].Status)
			},
		},
		{
			name: "not found",
			mock: func(repo *mocks.MockRepository) {
				repo.EXPECT().FindByID(gomock.Any(), "nonexistent").Return(nil, domain.ErrShipmentNotFound)
			},
			id: "nonexistent",
			check: func(t *testing.T, _ []domain.Event, err error) {
				assert.ErrorIs(t, err, domain.ErrShipmentNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockRepository(ctrl)
			tt.mock(repo)

			svc := appshipment.NewService(repo)
			events, err := svc.GetShipmentEvents(context.Background(), tt.id)
			tt.check(t, events, err)
		})
	}
}
