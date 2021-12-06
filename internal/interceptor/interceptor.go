package interceptor

import (
	"context"
	"sort"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
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
}

func NewInterceptor(log *zap.Logger, persist *repo.Repo) *Interceptor {
	return &Interceptor{
		log:     log,
		persist: persist,
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
		authors := inter.getInfoAuthorization(ctx)
		if len(authors) == 0 {
			return ctx, nil
		}

		if roles, ok := authors[fullMethod]; ok {
			if len(roles) == 0 {
				return ctx, nil
			}

			_, err := grpc_auth.AuthFromMD(ctx, common.TokenType)
			if err != nil {
				return nil, err
			}

			// _, err = common.VerifyToken(inter.log, token)
			// if err != nil {
			// 	return nil, err
			// }

			// // userID := cast.ToUint64(verifiedToken.StandardClaims.Subject)
			// // if token != "123456789" {
			// // 	return nil, status.Errorf(codes.PermissionDenied, "bad token")
			// // }

			userRole := "admin"
			if inter.searchRoles(roles, userRole) {
				return ctx, nil
			}
		}

		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}
}

func (inter *Interceptor) getInfoAuthorization(ctx context.Context) map[string][]string {
	authorize := map[string][]string{}

	permissions := inter.persist.FindAllPermissions()
	for _, per := range permissions {
		authorize[per.Slug] = inter.getPermissionRoles(ctx, per)
	}

	inter.log.Info("Permission",
		zap.Any("authorize", authorize),
	)

	return authorize
}

func (inter *Interceptor) getPermissionRoles(ctx context.Context, per *ent.Permission) (roles []string) {
	perRoles, _ := per.QueryRoles().Where(role.DeleteTimeIsNil()).All(ctx)
	for _, perRole := range perRoles {
		roles = append(roles, perRole.Slug)
	}

	return roles
}

func (inter *Interceptor) searchRoles(roles []string, userRole string) bool {
	sort.Strings(roles)
	i := sort.Search(len(roles), func(i int) bool { return roles[i] >= userRole })
	if i < len(roles) && roles[i] == userRole {
		return true
	}

	return false
}
