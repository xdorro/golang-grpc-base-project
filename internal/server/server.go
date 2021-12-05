package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/kucow/golang-grpc-base/internal/common"
	"github.com/kucow/golang-grpc-base/internal/service"
)

// func buildDummyAuthFunction(expectedScheme string, expectedToken string) func(ctx context.Context) (
// 	context.Context, error,
// ) {
// 	return func(ctx context.Context) (context.Context, error) {
// 		token, err := grpc_auth.AuthFromMD(ctx, expectedScheme)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if token != expectedToken {
// 			return nil, status.Errorf(codes.PermissionDenied, "buildDummyAuthFunction bad token")
// 		}
// 		return context.WithValue(ctx, "some_context_marker", "marker_exists"), nil
// 	}
// }
//
// type server struct {
// 	helloworld.UnimplementedGreeterServer
// }
//
// func NewServer(opts *common.Option) *server {
// 	srv := &server{}
//
// 	// Create a listener on TCP port
// 	lis, err := net.Listen("tcp", ":8080")
// 	if err != nil {
// 		log.Fatalln("Failed to listen:", err)
// 	}
//
// 	alwaysLoggingDeciderServer := func(
// 		ctx context.Context, fullMethodName string, servingObject interface{},
// 	) bool {
// 		return true
// 	}
//
// 	srvOpts := []grpc.ServerOption{
// 		grpc_middleware.WithStreamServerChain(
// 			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
// 			grpc_zap.StreamServerInterceptor(opts.Log),
// 			grpc_zap.PayloadStreamServerInterceptor(opts.Log, alwaysLoggingDeciderServer),
// 			grpc_auth.StreamServerInterceptor(buildDummyAuthFunction("bearer", "some_good_token")),
// 		),
// 		grpc_middleware.WithUnaryServerChain(
// 			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
// 			grpc_zap.UnaryServerInterceptor(opts.Log),
// 			grpc_zap.PayloadUnaryServerInterceptor(opts.Log, alwaysLoggingDeciderServer),
// 			grpc_auth.UnaryServerInterceptor(buildDummyAuthFunction("bearer", "some_good_token")),
// 		),
// 	}
//
// 	// Create a gRPC server object
// 	s := grpc.NewServer(srvOpts...)
// 	// Attach the Greeter service to the server
// 	helloworld.RegisterGreeterServer(s, srv)
// 	// Serve gRPC server
// 	log.Println("Serving gRPC on 0.0.0.0:8080")
// 	go func() {
// 		log.Fatalln(s.Serve(lis))
// 	}()
//
// 	// Create a client connection to the gRPC server we just started
// 	// This is where the gRPC-Gateway proxies the requests
// 	conn, err := grpc.DialContext(
// 		context.Background(),
// 		"0.0.0.0:8080",
// 		grpc.WithBlock(),
// 		grpc.WithInsecure(),
// 	)
// 	if err != nil {
// 		log.Fatalln("Failed to dial server:", err)
// 	}
//
// 	gwmux := runtime.NewServeMux()
// 	// Register Greeter
// 	err = helloworld.RegisterGreeterHandler(context.Background(), gwmux, conn)
// 	if err != nil {
// 		log.Fatalln("Failed to register gateway:", err)
// 	}
//
// 	gwServer := &http.Server{
// 		Addr:    ":8090",
// 		Handler: gwmux,
// 	}
//
// 	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
// 	log.Fatalln(gwServer.ListenAndServe())
// 	return srv
// }
//
// func (s *server) SayHello(_ context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
// 	return &helloworld.HelloReply{Message: in.Name + " world"}, nil
// }

// Server struct
type Server struct {
	ctx context.Context
	log *zap.Logger

	grpcServer *grpc.Server
	httpServer *http.Server
}

// NewServer create new server
func NewServer(opts *common.Option) (*Server, error) {
	srv := &Server{
		ctx: opts.Ctx,
		log: opts.Log,
	}

	grpcPort := viper.GetInt("GRPC_PORT")
	srv.log.Info(fmt.Sprintf("Serving gRPC on http://localhost:%d", grpcPort))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		return nil, fmt.Errorf("net.Listen(): %w", err)
	}

	svc := service.NewService(opts)

	go func() {
		if err = srv.createServer(listener, svc); err != nil {
			opts.Log.Fatal("createServer()", zap.Error(err))
		}
	}()

	go func() {
		if err = srv.listenClient(listener); err != nil {
			opts.Log.Fatal("listenClient()", zap.Error(err))
		}
	}()

	return srv, nil
}

func (srv *Server) Close() error {
	srv.grpcServer.GracefulStop()

	if err := srv.httpServer.Shutdown(srv.ctx); err != nil {
		srv.log.Error("srv.restServer.Shutdown()", zap.Error(err))
		return err
	}

	return nil
}
