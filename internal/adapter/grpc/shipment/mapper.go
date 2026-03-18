package shipment

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "github.com/canyouhearthemusic/shipment-service/gen/shipment/v1"
	domain "github.com/canyouhearthemusic/shipment-service/internal/domain/shipment"
)

func domainToProtoShipment(s *domain.Shipment) *v1.Shipment {
	return &v1.Shipment{
		Id:              s.ID,
		ReferenceNumber: s.ReferenceNumber,
		Origin:          s.Origin,
		Destination:     s.Destination,
		Status:          domainToProtoStatus(s.CurrentStatus),
		DriverName:      s.DriverName,
		UnitNumber:      s.UnitNumber,
		Amount:          s.Amount,
		DriverRevenue:   s.DriverRevenue,
		CreatedAt:       timestamppb.New(s.CreatedAt),
		UpdatedAt:       timestamppb.New(s.UpdatedAt),
	}
}

func domainToProtoEvent(e domain.Event) *v1.ShipmentEvent {
	return &v1.ShipmentEvent{
		Id:         e.ID,
		ShipmentId: e.ShipmentID,
		Status:     domainToProtoStatus(e.Status),
		OccurredAt: timestamppb.New(e.OccurredAt),
	}
}

func domainToProtoStatus(s domain.Status) v1.ShipmentStatus {
	switch s {
	case domain.StatusPending:
		return v1.ShipmentStatus_SHIPMENT_STATUS_PENDING

	case domain.StatusPickedUp:
		return v1.ShipmentStatus_SHIPMENT_STATUS_PICKED_UP

	case domain.StatusInTransit:
		return v1.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT

	case domain.StatusDelivered:
		return v1.ShipmentStatus_SHIPMENT_STATUS_DELIVERED

	case domain.StatusCancelled:
		return v1.ShipmentStatus_SHIPMENT_STATUS_CANCELLED

	default:
		return v1.ShipmentStatus_SHIPMENT_STATUS_UNSPECIFIED
	}
}

func protoToDomainStatus(s v1.ShipmentStatus) domain.Status {
	switch s {
	case v1.ShipmentStatus_SHIPMENT_STATUS_PENDING:
		return domain.StatusPending

	case v1.ShipmentStatus_SHIPMENT_STATUS_PICKED_UP:
		return domain.StatusPickedUp

	case v1.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT:
		return domain.StatusInTransit

	case v1.ShipmentStatus_SHIPMENT_STATUS_DELIVERED:
		return domain.StatusDelivered

	case v1.ShipmentStatus_SHIPMENT_STATUS_CANCELLED:
		return domain.StatusCancelled

	default:
		return ""
	}
}
