package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/go-redis/redis/v8"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/spf13/cast"
	"github.com/vk-rv/pvx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/xdorro/golang-grpc-base-project/api/ent"
	"github.com/xdorro/golang-grpc-base-project/api/ent/role"
	"github.com/xdorro/golang-grpc-base-project/internal/common"
	"github.com/xdorro/golang-grpc-base-project/internal/handler"
	"github.com/xdorro/golang-grpc-base-project/internal/repo"
)

type Interceptor struct {
	log     *zap.Logger
	redis   redis.UniversalClient
	repo    *repo.Repo
	handler *handler.Handler
	mutex   *sync.Mutex
}

func NewInterceptor(
	log *zap.Logger, repo *repo.Repo, redis redis.UniversalClient, handler *handler.Handler,
) *Interceptor {
	return &Interceptor{
		log:     log,
		repo:    repo,
		redis:   redis,
		handler: handler,
		mutex:   &sync.Mutex{},
	}
}

// AuthInterceptorStream create auth Interceptor stream
func (interceptor *Interceptor) AuthInterceptorStream() grpc.StreamServerInterceptor {
	return func(
		grpcSrv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler,
	) error {
		var newCtx context.Context
		var err error
		if overrideSrv, ok := grpcSrv.(grpc_auth.ServiceAuthFuncOverride); ok {
			newCtx, err = overrideSrv.AuthFuncOverride(stream.Context(), info.FullMethod)
		} else {
			authFunc := interceptor.authInterceptor(info.FullMethod)
			newCtx, err = authFunc(stream.Context())
		}
		if err != nil {
			return err
		}
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		return handler(grpcSrv, wrapped)
	}
}

// AuthInterceptorUnary create auth Interceptor unary
func (interceptor *Interceptor) AuthInterceptorUnary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (interface{}, error) {
		var newCtx context.Context
		var err error
		if overrideSrv, ok := info.Server.(grpc_auth.ServiceAuthFuncOverride); ok {
			newCtx, err = overrideSrv.AuthFuncOverride(ctx, info.FullMethod)
		} else {
			authFunc := interceptor.authInterceptor(info.FullMethod)
			newCtx, err = authFunc(ctx)
		}
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

// authInterceptor handler interceptor
func (interceptor *Interceptor) authInterceptor(fullMethod string) grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		authorize := interceptor.getInfoAuthorization(ctx)
		if len(authorize) == 0 {
			return ctx, nil
		}

		if roles, ok := authorize[fullMethod]; ok {
			if len(roles) == 0 {
				return ctx, nil
			}

			accessToken, err := grpc_auth.AuthFromMD(ctx, common.TokenType)
			if err != nil {
				return nil, err
			}

			var claims *pvx.RegisteredClaims
			claims, err = interceptor.handler.DecryptToken(accessToken)
			if err != nil {
				return nil, status.New(codes.InvalidArgument, err.Error()).Err()
			}

			interceptor.log.Info("svc.handler.DecryptToken()",
				zap.Any("claims", claims),
			)

			tokenKey := fmt.Sprintf(common.UserSessionKey, claims.Subject, claims.TokenID)
			if err = interceptor.handler.ExistRefreshToken(tokenKey); err != nil {
				return nil, err
			}

			var u *ent.User
			u, err = interceptor.repo.FindUserByID(cast.ToUint64(claims.Subject))
			if err != nil {
				return nil, common.EmailNotExist.Err()
			}

			// add key-value pairs of metadata to context
			ctx = metadata.NewOutgoingContext(
				ctx,
				metadata.Pairs(common.CtxUserID, claims.Subject),
			)

			var userRole *ent.Role
			userRole, err = u.QueryRoles().Where(role.DeleteTimeIsNil()).First(ctx)
			if err != nil {
				return ctx, err
			}

			if hasAccessTo(roles, userRole) {
				return ctx, nil
			}
		}

		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}
}

// getInfoAuthorization get info authorization
func (interceptor *Interceptor) getInfoAuthorization(ctx context.Context) map[string][]string {
	authorize := make(map[string][]string)

	if val := interceptor.redis.Get(ctx, common.KeyServiceRoles).Val(); val != "" {
		authorize = cast.ToStringMapStringSlice(val)
		return authorize
	}

	permissions := interceptor.repo.FindAllPermissionsWithRoles()
	for _, per := range permissions {
		authorize[per.Slug] = getPermissionRoles(per)
	}

	data, _ := json.Marshal(authorize)
	if err := interceptor.redis.Set(ctx, common.KeyServiceRoles, data, common.KeyServiceRolesExpire).Err(); err != nil {
		interceptor.log.Error("redis.Set()", zap.Error(err))
	}

	return authorize
}

// getPermissionRoles get permission roles
func getPermissionRoles(per *ent.Permission) []string {
	roles := make([]string, 0)

	perRoles := per.Edges.Roles
	for _, perRole := range perRoles {
		roles = append(roles, perRole.Slug)
	}

	return roles
}

// hasAccessTo check has access
func hasAccessTo(roles []string, role *ent.Role) bool {
	if role.FullAccess {
		return true
	}

	for _, r := range roles {
		if strings.EqualFold(role.Slug, r) {
			return true
		}
	}

	return false
}
