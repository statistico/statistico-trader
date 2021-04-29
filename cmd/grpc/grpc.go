package main

import (
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-trader/internal/trader/bootstrap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

func main() {
	app := bootstrap.BuildContainer(bootstrap.BuildConfig())
	logger := app.Logger

	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
	}

	opts := grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionIdle:5*time.Minute})
	grpcServer := grpc.NewServer(opts)
	statistico.RegisterStrategyServiceServer(grpcServer, app.GrpcStrategyService())
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatalf("Failed to serve: %v", err)
	}
}
