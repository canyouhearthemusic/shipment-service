# Shipment Service

gRPC microservice for tracking shipments and their status changes throughout the delivery lifecycle. Built with Go.

## Running the service

You'll need Go 1.25+ installed.

```bash
make run
```

This starts the gRPC server on port 50051 by default. You can change it with the `GRPC_PORT` env variable.

With Docker:

```bash
docker compose up --build
```

To regenerate protobuf code (requires `protoc`, `protoc-gen-go`, `protoc-gen-go-grpc`):

```bash
make proto
```

To regenerate mocks:

```bash
make generate
```

You can poke at the running server with [grpcurl](https://github.com/fullstorydev/grpcurl):

```bash
grpcurl -plaintext -d '{
  "reference_number": "REF-001",
  "origin": "NYC",
  "destination": "LA",
  "driver_name": "John",
  "unit_number": "TRK-1",
  "amount": 1500,
  "driver_revenue": 300
}' localhost:50051 shipment.v1.ShipmentService/CreateShipment

grpcurl -plaintext -d '{"id": "<shipment-id>"}' \
  localhost:50051 shipment.v1.ShipmentService/GetShipment
```

## Running the tests

```bash
make test
```

This runs all tests with the race detector enabled. Tests cover two layers:

- **Domain** — shipment creation, the full status transition state machine (valid and invalid), event recording
- **Application** — service orchestration with a mocked repository (gomock), verifying correct delegation and error propagation
