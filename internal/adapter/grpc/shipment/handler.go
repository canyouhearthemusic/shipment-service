package shipment

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	v1 "github.com/canyouhearthemusic/shipment-service/gen/shipment/v1"
	application "github.com/canyouhearthemusic/shipment-service/internal/application/shipment"
	domain "github.com/canyouhearthemusic/shipment-service/internal/domain/shipment"
)

func Register(srv *grpc.Server, svc application.Service) {
	v1.RegisterShipmentServiceServer(srv, newHandler(svc))
}

type handler struct {
	v1.UnimplementedShipmentServiceServer
	svc application.Service
}

func newHandler(svc application.Service) *handler {
	return &handler{svc: svc}
}

func (h *handler) CreateShipment(ctx context.Context, req *v1.CreateShipmentRequest) (*v1.CreateShipmentResponse, error) {
	s, err := h.svc.CreateShipment(ctx, application.CreateShipmentRequest{
		ReferenceNumber: req.ReferenceNumber,
		Origin:          req.Origin,
		Destination:     req.Destination,
		DriverName:      req.DriverName,
		UnitNumber:      req.UnitNumber,
		Amount:          req.Amount,
		DriverRevenue:   req.DriverRevenue,
	})

	if err != nil {
		return nil, toGRPCError(err)
	}

	return &v1.CreateShipmentResponse{Shipment: domainToProtoShipment(s)}, nil
}

func (h *handler) GetShipment(ctx context.Context, req *v1.GetShipmentRequest) (*v1.GetShipmentResponse, error) {
	s, err := h.svc.GetShipment(ctx, req.Id)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &v1.GetShipmentResponse{Shipment: domainToProtoShipment(s)}, nil
}

func (h *handler) AddShipmentEvent(ctx context.Context, req *v1.AddShipmentEventRequest) (*v1.AddShipmentEventResponse, error) {
	domainStatus := protoToDomainStatus(req.Status)
	if domainStatus == "" {
		return nil, status.Error(codes.InvalidArgument, "unknown shipment status")
	}

	s, err := h.svc.AddShipmentEvent(ctx, req.Id, domainStatus)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &v1.AddShipmentEventResponse{Shipment: domainToProtoShipment(s)}, nil
}

func (h *handler) GetShipmentEvents(ctx context.Context, req *v1.GetShipmentEventsRequest) (*v1.GetShipmentEventsResponse, error) {
	events, err := h.svc.GetShipmentEvents(ctx, req.Id)
	if err != nil {
		return nil, toGRPCError(err)
	}

	protoEvents := make([]*v1.ShipmentEvent, len(events))
	for i, e := range events {
		protoEvents[i] = domainToProtoEvent(e)
	}

	return &v1.GetShipmentEventsResponse{Events: protoEvents}, nil
}

func toGRPCError(err error) error {
	switch {
	case errors.Is(err, domain.ErrShipmentNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, domain.ErrInvalidTransition):
		return status.Error(codes.FailedPrecondition, err.Error())

	case errors.Is(err, domain.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())

	default:
		return status.Error(codes.Internal, err.Error())
	}
}
