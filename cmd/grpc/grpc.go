package main

import (
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-trader/internal/trader/bootstrap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
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

	lis, err := net.Listen("tcp", ":7777")

	if err != nil {
		app.Logger.Fatalf("Failed to listen: %v", err)
	}

	opts := grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionIdle:5*time.Minute})
	grpcServer := grpc.NewServer(opts)
	statistico.RegisterStrategyServiceServer(grpcServer, app.GrpcStrategyService())
	reflection.Register(grpcServer)

	go func() {
		log.Fatal(grpcServer.Serve(lis))
	}()

	grpcWebServer := grpcweb.WrapServer(
		grpcServer,
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
		grpcweb.WithOriginFunc(func(origin string) bool {
			return true
		}),
	)

	multiplex := grpcMultiplexer{grpcWebServer}

	srv := &http.Server{
		Handler:      h2c.NewHandler(multiplex.Handler(), &http2.Server{}),
		Addr:         ":50051",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  600 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

type grpcMultiplexer struct {
	*grpcweb.WrappedGrpcServer
}

func (m *grpcMultiplexer) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request URI %s", r.RequestURI)
		log.Printf("Request Method %s", r.Method)
		if r.Method == http.MethodOptions {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Accept, Origin, Authorization, X-User-Agent,X-Grpc-Web, Authorization, Keep-Alive")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
			return
		}

		if m.IsGrpcWebRequest(r) || r.ProtoMajor == 2 {
			m.ServeHTTP(w, r)
			return
		}

		return
	})
}
