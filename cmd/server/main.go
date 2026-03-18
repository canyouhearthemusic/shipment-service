package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/canyouhearthemusic/shipment-service/config"
	"github.com/canyouhearthemusic/shipment-service/internal/adapter/grpc"
	shipmentgrpc "github.com/canyouhearthemusic/shipment-service/internal/adapter/grpc/shipment"
	"github.com/canyouhearthemusic/shipment-service/internal/adapter/persistence/memory"
	shipmentapp "github.com/canyouhearthemusic/shipment-service/internal/application/shipment"
)

func main() {
	cfg := config.Load()

	shipmentRepo := memory.NewShipmentRepository()
	shipmentSvc := shipmentapp.NewService(shipmentRepo)

	srv := grpc.NewServer()
	shipmentgrpc.Register(srv, shipmentSvc)

	addr := fmt.Sprintf(":%s", cfg.GRPC.Port)
	log.Printf("gRPC server listening on %s", addr)

	go func() {
		if err := grpc.ListenAndServe(srv, addr); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down gracefully...")
	srv.GracefulStop()
}
