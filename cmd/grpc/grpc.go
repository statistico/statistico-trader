package main

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
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

	auth := app.TokenAuthoriser()

	opts := grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionIdle:5*time.Minute})

	server := grpc.NewServer(
		opts,
		grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(auth.Authorise)),
		grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(auth.Authorise)),
	)

	statistico.RegisterStrategyServiceServer(server, app.GrpcStrategyService())

	reflection.Register(server)

	if err := server.Serve(lis); err != nil {
		logger.Fatalf("Failed to serve: %v", err)
	}
}
