package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/xdorro/golang-grpc-base-project/api/ent"
	auth_proto "github.com/xdorro/golang-grpc-base-project/api/proto/v1/auth"
	permission_proto "github.com/xdorro/golang-grpc-base-project/api/proto/v1/permission"
	role_proto "github.com/xdorro/golang-grpc-base-project/api/proto/v1/role"
	user_proto "github.com/xdorro/golang-grpc-base-project/api/proto/v1/user"
	"github.com/xdorro/golang-grpc-base-project/internal/common"
	"github.com/xdorro/golang-grpc-base-project/internal/repo"
	"github.com/xdorro/golang-grpc-base-project/internal/service"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewServer, NewGRPCServer, NewHTTPServer, NewInterceptor)

// Server struct
type Server struct {
	ctx   context.Context
	redis redis.UniversalClient

	log   *zap.Logger
	repo  *repo.Repo
	grpc  *grpc.Server
	mutex *sync.Mutex
}

func NewServer(
	ctx context.Context, log *zap.Logger, repo *repo.Repo, redis redis.UniversalClient, grpc *grpc.Server,
	mux *runtime.ServeMux, _ *service.Service,
) (*Server, error) {
	srv := &Server{
		ctx:   ctx,
		log:   log,
		repo:  repo,
		redis: redis,
		grpc:  grpc,
		mutex: &sync.Mutex{},
	}

	grpcPort := fmt.Sprintf(":%d", viper.GetInt("GRPC_PORT"))
	log.Info(fmt.Sprintf("Serving gRPC on http://localhost%s", grpcPort))

	listenGRPC, err := net.Listen("tcp", grpcPort)
	if err != nil {
		return nil, fmt.Errorf("net.Listen(): %w", err)
	}

	httpPort := fmt.Sprintf(":%d", viper.GetInt("HTTP_PORT"))
	log.Info(fmt.Sprintf("Serving gRPC-Gateway on http://localhost%s", httpPort))

	listenHTTP, err := net.Listen("tcp", httpPort)
	if err != nil {
		return nil, fmt.Errorf("net.Listen(): %w", err)
	}

	// get UserService Info
	if viper.GetBool("SEEDER_SERVICE") {
		go srv.getServiceInfo()
	}

	go func() {
		if err = grpc.Serve(listenGRPC); err != nil {
			srv.log.Fatal("grpc.Serve()", zap.Error(err))
		}
	}()

	go func() {
		if err = srv.registerServiceHandlers(grpcPort, mux); err != nil {
			log.Fatal("srv.registerServiceHandlers()", zap.Error(err))
		}

		if err = http.Serve(listenHTTP, mux); err != nil {
			log.Fatal("http.Serve()", zap.Error(err))
		}
	}()

	return srv, nil
}

func (srv *Server) Close() error {
	if srv.grpc != nil {
		srv.grpc.GracefulStop()
	}

	if srv.repo != nil {
		if err := srv.repo.Close(); err != nil {
			return err
		}
	}

	if srv.redis != nil {
		if err := srv.redis.Close(); err != nil {
			return err
		}
	}

	return nil
}

// getServiceInfo returns service info
func (srv *Server) getServiceInfo() {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	bulk := make([]*ent.PermissionCreate, 0)
	for name, val := range srv.grpc.GetServiceInfo() {
		if len(val.Methods) == 0 {
			return
		}

		for _, info := range val.Methods {
			slug := fmt.Sprintf("/%s/%s", name, info.Name)

			if !srv.repo.ExistPermissionBySlug(slug) {
				srv.log.Info("GetServiceInfo",
					zap.Any("Name", info.Name),
					zap.Any("Slug", slug),
				)

				bulk = append(bulk, srv.repo.Permission.
					Create().
					SetName(info.Name).
					SetSlug(slug).
					SetStatus(1),
				)
			}
		}
	}

	if len(bulk) > 0 {
		if err := srv.repo.CreatePermissionBulk(bulk); err != nil {
			srv.log.Error("persist.CreatePermissionBulk()", zap.Error(err))
		}

		if err := srv.redis.Del(srv.ctx, common.KeyServiceRoles).Err(); err != nil {
			srv.log.Error("redis.Del()", zap.Error(err))
		}
	}
}

func (srv *Server) registerServiceHandlers(grpcPort string, mux *runtime.ServeMux) error {
	conn, err := grpc.DialContext(
		srv.ctx,
		grpcPort,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		srv.log.Fatal("Failed to dial Server:", zap.Error(err))
	}

	// Register AuthService Handler
	if err = auth_proto.RegisterAuthServiceHandler(srv.ctx, mux, conn); err != nil {
		return fmt.Errorf("proto.RegisterAuthServiceHandler(): %w", err)
	}

	// Register UserService Handler
	if err = user_proto.RegisterUserServiceHandler(srv.ctx, mux, conn); err != nil {
		return fmt.Errorf("proto.RegisterUserServiceHandler(): %w", err)
	}

	// Register RoleService Handler
	if err = role_proto.RegisterRoleServiceHandler(srv.ctx, mux, conn); err != nil {
		return fmt.Errorf("proto.RegisterRoleServiceHandler(): %w", err)
	}

	// Register PermissionService Handler
	if err = permission_proto.RegisterPermissionServiceHandler(srv.ctx, mux, conn); err != nil {
		return fmt.Errorf("proto.RegisterPermissionServiceHandler(): %w", err)
	}

	return nil
}
