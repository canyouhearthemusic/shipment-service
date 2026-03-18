package shipment

type CreateShipmentRequest struct {
	ReferenceNumber string
	Origin          string
	Destination     string
	DriverName      string
	UnitNumber      string
	Amount          float64
	DriverRevenue   float64
}
