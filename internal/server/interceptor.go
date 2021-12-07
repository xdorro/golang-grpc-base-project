package server

import (
	"context"
	"encoding/json"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xdorro/golang-grpc-base-project/ent"
	"github.com/xdorro/golang-grpc-base-project/ent/role"
	"github.com/xdorro/golang-grpc-base-project/internal/common"
)

type contextKey string

const (
	ctxUserID contextKey = "userID"
)

// AuthInterceptorStream create auth Interceptor stream
func (srv *server) AuthInterceptorStream() grpc.StreamServerInterceptor {
	return func(
		grpcSrv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler,
	) error {
		var newCtx context.Context
		var err error
		if overrideSrv, ok := grpcSrv.(grpc_auth.ServiceAuthFuncOverride); ok {
			newCtx, err = overrideSrv.AuthFuncOverride(stream.Context(), info.FullMethod)
		} else {
			authFunc := srv.authInterceptor(info.FullMethod)
			newCtx, err = authFunc(stream.Context())
		}
		if err != nil {
			return err
		}
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		return handler(srv, wrapped)
	}
}

// AuthInterceptorUnary create auth Interceptor unary
func (srv *server) AuthInterceptorUnary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (interface{}, error) {
		var newCtx context.Context
		var err error
		if overrideSrv, ok := info.Server.(grpc_auth.ServiceAuthFuncOverride); ok {
			newCtx, err = overrideSrv.AuthFuncOverride(ctx, info.FullMethod)
		} else {
			authFunc := srv.authInterceptor(info.FullMethod)
			newCtx, err = authFunc(ctx)
		}
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

// authInterceptor handler interceptor
func (srv *server) authInterceptor(fullMethod string) grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		authorize := srv.getInfoAuthorization(ctx)

		if len(authorize) == 0 {
			return ctx, nil
		}

		if roles, ok := authorize[fullMethod]; ok {
			if len(roles) == 0 {
				return ctx, nil
			}

			token, err := grpc_auth.AuthFromMD(ctx, common.TokenType)
			if err != nil {
				return nil, err
			}

			verifiedToken, err := common.VerifyToken(srv.Log, token)
			if err != nil {
				return nil, err
			}

			claims := verifiedToken.StandardClaims
			user, err := srv.Persist.FindUserByID(cast.ToUint64(claims.Subject))
			if err != nil {
				return nil, common.UserNotExist.Err()
			}

			userRoles, _ := user.QueryRoles().Where(role.DeleteTimeIsNil()).All(ctx)

			if srv.hasAccessTo(roles, userRoles) {
				return context.WithValue(ctx, ctxUserID, user.ID), nil
			}
		}

		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}
}

// getInfoAuthorization get info authorization
func (srv *server) getInfoAuthorization(ctx context.Context) map[string][]string {
	authorize := make(map[string][]string)

	val := srv.Redis.Get(srv.Redis.Context(), common.ServiceRoles).Val()
	if val != "" {
		authorize = cast.ToStringMapStringSlice(val)
		return authorize
	}

	permissions := srv.Persist.FindAllPermissions()
	for _, per := range permissions {
		authorize[per.Slug] = srv.getPermissionRoles(ctx, per)
	}

	data, _ := json.Marshal(authorize)
	err := srv.Redis.Set(srv.Redis.Context(), common.ServiceRoles, data, -1).Err()
	if err != nil {
		srv.Log.Error("redis.Set()", zap.Error(err))
	}

	return authorize
}

// getPermissionRoles get permission roles
func (srv *server) getPermissionRoles(ctx context.Context, per *ent.Permission) []string {
	roles := make([]string, 0)

	perRoles, _ := per.QueryRoles().Where(role.DeleteTimeIsNil()).All(ctx)
	for _, perRole := range perRoles {
		roles = append(roles, perRole.Slug)
	}

	return roles
}

// hasAccessTo check has access
func (srv *server) hasAccessTo(roles []string, userRoles []*ent.Role) bool {
	for _, ur := range userRoles {
		if ur.FullAccess {
			return true
		}

		for _, r := range roles {
			if strings.EqualFold(ur.Slug, r) {
				return true
			}
		}
	}

	return false
}
