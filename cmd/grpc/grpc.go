package main

import (
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-trader/internal/trader/bootstrap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

func main() {
	app := bootstrap.BuildContainer(bootstrap.BuildConfig())

	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		app.Logger.Fatalf("Failed to listen: %v", err)
	}

	opts := grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionIdle:5*time.Minute})

	server := grpc.NewServer(opts)

	statistico.RegisterStrategyServiceServer(server, app.GrpcStrategyService())

	reflection.Register(server)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
