package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/xdorro/golang-grpc-base-project/internal/log"
	"github.com/xdorro/golang-grpc-base-project/internal/service"
)

// Server is server struct.
type Server struct {
	ctx        context.Context
	grpcServer *grpc.Server
	httpServer *runtime.ServeMux
	mu         *sync.RWMutex
}

// NewServer creates a new server.
func NewServer(ctx context.Context, service service.Service) IServer {
	s := &Server{
		ctx: ctx,
		mu:  &sync.RWMutex{},
	}

	// pprof debug mode
	if viper.GetBool("APP_DEBUG") {
		go func() {
			if err := http.ListenAndServe("localhost:6060", nil); err != nil {
				log.Panicf("pprof server error: %v", err)
			}
		}()
	}

	cert := viper.GetString("APP_CERT")
	key := viper.GetString("APP_KEY")
	tlsCredentials, err := loadTLSCredentials(cert, key)
	if err != nil {
		log.Panicf("cannot load TLS credentials: %v", err)
	}

	appPort := fmt.Sprintf(":%d", viper.GetInt("APP_PORT"))
	log.Infof("Serving on https://localhost%s", appPort)
	handler := s.httpGrpcRouter(tlsCredentials, appPort, service)

	go func() {
		if err = http.ListenAndServeTLS(appPort, cert, key, handler); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Panicf("error serving: %v", err)
		}
	}()

	return s
}

// Close closes the server.
func (s *Server) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Infof("Closing server...")
	s.grpcServer.GracefulStop()

	return nil
}

// httpGrpcRouter is http grpc router.
func (s *Server) httpGrpcRouter(
	tlsCredentials credentials.TransportCredentials, appPort string, service service.Service,
) http.Handler {
	s.newHTTPServer(tlsCredentials, appPort)
	s.newGRPCServer(tlsCredentials, service)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
			s.grpcServer.ServeHTTP(w, r)
			return
		}

		// middleware that adds CORS headers to the response.
		h := w.Header()
		h.Set("Access-Control-Allow-Origin", "http://localhost:3000")
		h.Set("Access-Control-Allow-Credentials", "true")

		if strings.EqualFold(r.Method, http.MethodOptions) {
			h.Set("Access-Control-Methods", "POST, PUT, PATCH, DELETE")
			h.Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,Content-Type")
			h.Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		s.httpServer.ServeHTTP(w, r)
	})
}

// loadTLSCredentials loads TLS credentials from the configuration
func loadTLSCredentials(cert, key string) (
	credentials.TransportCredentials, error,
) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates:       []tls.Certificate{serverCert},
		ClientAuth:         tls.RequireAndVerifyClientCert,
		InsecureSkipVerify: true,
	}

	return credentials.NewTLS(config), nil
}
