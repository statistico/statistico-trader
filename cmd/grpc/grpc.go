package main

import (
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-trader/internal/trader/bootstrap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	app := bootstrap.BuildContainer(bootstrap.BuildConfig())

	lis, err := net.Listen("tcp", ":50052")

	if err != nil {
		app.Logger.Fatalf("Failed to listen: %v", err)
	}

	opts := grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionIdle:5*time.Minute})

	server := grpc.NewServer(opts)

	statistico.RegisterStrategyServiceServer(server, app.GrpcStrategyService())

	reflection.Register(server)

	go func() {
		log.Fatal(server.Serve(lis))
	}()

	grpcWebServer := grpcweb.WrapServer(
		server,
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
		grpcweb.WithOriginFunc(func(origin string) bool {
			return origin == "http://localhost:3000" || origin == "http://localhost:3000/"
		}),
	)

	srv := &http.Server{
		Addr:         "localhost:50051",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	srv.Handler = http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if grpcWebServer.IsGrpcWebRequest(req) {
			grpcWebServer.ServeHTTP(resp, req)
		}
		// Fall back to other servers.
		http.DefaultServeMux.ServeHTTP(resp, req)
	})

	log.Fatal(srv.ListenAndServe())
}
