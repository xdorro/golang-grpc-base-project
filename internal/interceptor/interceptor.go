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

	"github.com/kucow/golang-grpc-base-project/internal/common"
	"github.com/kucow/golang-grpc-base-project/internal/repo"
	"github.com/kucow/golang-grpc-base-project/pkg/ent"
	"github.com/kucow/golang-grpc-base-project/pkg/ent/role"
)

type Interceptor struct {
	log     *zap.Logger
	persist repo.Persist
	redis   redis.UniversalClient
}

type contextKey string

const (
	ctxUserID contextKey = "userID"
)

func NewInterceptor(log *zap.Logger, redis redis.UniversalClient, persist *repo.Repo) *Interceptor {
	return &Interceptor{
		log:     log,
		persist: persist,
		redis:   redis,
	}
}

func (inter *Interceptor) AuthInterceptorStream() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
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

func (inter *Interceptor) AuthInterceptorUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
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

func (inter *Interceptor) authInterceptor(fullMethod string) grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		authorize := inter.getInfoAuthorization(ctx)
		inter.log.Info("Permission", zap.Any("authorize", authorize))

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
			if inter.hasAccessTo(roles, claims.Audience) {
				userID := cast.ToUint64(claims.Subject)
				return context.WithValue(ctx, ctxUserID, userID), nil
			}
		}

		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}
}

func (inter *Interceptor) getInfoAuthorization(ctx context.Context) map[string][]string {
	authorize := map[string][]string{}
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

func (inter *Interceptor) getPermissionRoles(ctx context.Context, per *ent.Permission) (roles []string) {
	perRoles, _ := per.QueryRoles().Where(role.DeleteTimeIsNil()).All(ctx)
	for _, perRole := range perRoles {
		roles = append(roles, perRole.Slug)
	}

	return roles
}

func (inter *Interceptor) hasAccessTo(roles, userRoles []string) bool {
	for _, ur := range userRoles {
		if strings.EqualFold(ur, "admin") {
			return true
		}

		for _, r := range roles {
			if strings.EqualFold(ur, r) {
				return true
			}
		}
	}

	return false
}
