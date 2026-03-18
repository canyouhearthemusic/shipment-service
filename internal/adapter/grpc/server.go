package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"runtime/debug"

	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func NewServer(opts ...googlegrpc.ServerOption) *googlegrpc.Server {
	opts = append([]googlegrpc.ServerOption{
		googlegrpc.ChainUnaryInterceptor(
			unaryLoggingInterceptor,
			unaryRecoveryInterceptor,
		),
	}, opts...)

	srv := googlegrpc.NewServer(opts...)
	reflection.Register(srv)

	return srv
}

func ListenAndServe(srv *googlegrpc.Server, addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", addr, err)
	}

	return srv.Serve(lis)
}

func unaryLoggingInterceptor(ctx context.Context, req any, info *googlegrpc.UnaryServerInfo, handler googlegrpc.UnaryHandler) (any, error) {
	log.Printf("gRPC call: %s", info.FullMethod)
	resp, err := handler(ctx, req)
	if err != nil {
		log.Printf("gRPC error: %s: %v", info.FullMethod, err)
	}

	return resp, err
}

func unaryRecoveryInterceptor(ctx context.Context, req any, info *googlegrpc.UnaryServerInfo, handler googlegrpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic recovered in %s: %v\n%s", info.FullMethod, r, debug.Stack())
			err = status.Errorf(codes.Internal, "internal server error")
		}
	}()

	return handler(ctx, req)
}
