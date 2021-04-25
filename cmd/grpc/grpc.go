package main

import (
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-trader/internal/trader/bootstrap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"log"
	"net/http"
	"time"
)

func main() {
	app := bootstrap.BuildContainer(bootstrap.BuildConfig())
	logger := app.Logger

	opts := grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionIdle:5*time.Minute})
	server := grpc.NewServer(opts)

	statistico.RegisterStrategyServiceServer(server, app.GrpcStrategyService())
	reflection.Register(server)

	grpcWebServer := grpcweb.WrapServer(server,
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
		grpcweb.WithOriginFunc(func(origin string) bool {
			return origin == "https//localhost:3000" || origin == "http://localhost:3000"
		}),
	)

	srv := &http.Server{
		Handler:      http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			logger.Errorf("Request method %s", req.Method)
			if grpcWebServer.IsGrpcWebRequest(req) {
				logger.Error("Request is gRPC")
				grpcWebServer.ServeHTTP(resp, req)
				return
			}

			logger.Error("Request is NOT gRPC")
			// Fall back to other servers.
			http.DefaultServeMux.ServeHTTP(resp, req)
		}),
		Addr:         "localhost:50051",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  300 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Errorf("Listen and serve error %+v", err)
		log.Fatalf("Failed to serve: %v", err)
	}
}
