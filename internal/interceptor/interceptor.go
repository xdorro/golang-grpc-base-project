package interceptor

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/go-redis/redis/v8"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xdorro/golang-grpc-base-project/internal/common"
	"github.com/xdorro/golang-grpc-base-project/internal/repo"
	"github.com/xdorro/golang-grpc-base-project/pkg/ent"
	"github.com/xdorro/golang-grpc-base-project/pkg/ent/role"
)

// Interceptor struct
type Interceptor struct {
	log     *zap.Logger
	persist repo.Persist
	redis   redis.UniversalClient
}

type contextKey string

const (
	ctxUserID contextKey = "userID"
)

// NewInterceptor create new interceptor
func NewInterceptor(log *zap.Logger, redis redis.UniversalClient, persist *repo.Repo) *Interceptor {
	return &Interceptor{
		log:     log,
		persist: persist,
		redis:   redis,
	}
}

// AuthInterceptorStream create auth Interceptor stream
func (inter *Interceptor) AuthInterceptorStream() grpc.StreamServerInterceptor {
	return func(
		srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler,
	) error {
		var newCtx context.Context
		var err error
		if overrideSrv, ok := srv.(grpc_auth.ServiceAuthFuncOverride); ok {
			newCtx, err = overrideSrv.AuthFuncOverride(stream.Context(), info.FullMethod)
		} else {
			authFunc := inter.authInterceptor(info.FullMethod)
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
func (inter *Interceptor) AuthInterceptorUnary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (interface{}, error) {
		var newCtx context.Context
		var err error
		if overrideSrv, ok := info.Server.(grpc_auth.ServiceAuthFuncOverride); ok {
			newCtx, err = overrideSrv.AuthFuncOverride(ctx, info.FullMethod)
		} else {
			authFunc := inter.authInterceptor(info.FullMethod)
			newCtx, err = authFunc(ctx)
		}
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

// authInterceptor handler interceptor
func (inter *Interceptor) authInterceptor(fullMethod string) grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		authorize := inter.getInfoAuthorization(ctx)

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

			verifiedToken, err := common.VerifyToken(inter.log, token)
			if err != nil {
				return nil, err
			}

			claims := verifiedToken.StandardClaims
			user, err := inter.persist.FindUserByID(cast.ToUint64(claims.Subject))
			if err != nil {
				return nil, common.UserNotExist.Err()
			}

			userRoles, _ := user.QueryRoles().Where(role.DeleteTimeIsNil()).All(ctx)

			if inter.hasAccessTo(roles, userRoles) {
				return context.WithValue(ctx, ctxUserID, user.ID), nil
			}
		}

		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}
}

// getInfoAuthorization get info authorization
func (inter *Interceptor) getInfoAuthorization(ctx context.Context) map[string][]string {
	authorize := make(map[string][]string)

	val := inter.redis.Get(inter.redis.Context(), common.ServiceRoles).Val()
	if val != "" {
		authorize = cast.ToStringMapStringSlice(val)
		return authorize
	}

	permissions := inter.persist.FindAllPermissions()
	for _, per := range permissions {
		authorize[per.Slug] = inter.getPermissionRoles(ctx, per)
	}

	data, _ := json.Marshal(authorize)
	err := inter.redis.Set(inter.redis.Context(), common.ServiceRoles, data, -1).Err()
	if err != nil {
		inter.log.Error("redis.Set()", zap.Error(err))
	}

	return authorize
}

// getPermissionRoles get permission roles
func (inter *Interceptor) getPermissionRoles(ctx context.Context, per *ent.Permission) []string {
	roles := make([]string, 0)

	perRoles, _ := per.QueryRoles().Where(role.DeleteTimeIsNil()).All(ctx)
	for _, perRole := range perRoles {
		roles = append(roles, perRole.Slug)
	}

	return roles
}

// hasAccessTo check has access
func (inter *Interceptor) hasAccessTo(roles []string, userRoles []*ent.Role) bool {
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
