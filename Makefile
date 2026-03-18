PROTO_SRC  := proto/shipment/v1/shipment.proto
GEN_OUT    := .
MODULE     := github.com/canyouhearthemusic/shipment-service

.PHONY: proto build test lint run generate

proto:
	protoc \
		--proto_path=proto \
		--go_out=$(GEN_OUT) --go_opt=module=$(MODULE) \
		--go-grpc_out=$(GEN_OUT) --go-grpc_opt=module=$(MODULE) \
		$(PROTO_SRC)

build:
	go build ./...

test:
	go test -race ./...

lint:
	golangci-lint run

generate:
	go generate ./...

run:
	go run ./cmd/server
