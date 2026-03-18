package shipment

import "errors"

var (
	ErrInvalidTransition = errors.New("invalid status transition")
	ErrShipmentNotFound  = errors.New("shipment not found")
	ErrInvalidInput      = errors.New("invalid input")
)
