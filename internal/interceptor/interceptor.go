package interceptor

import (
	"context"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"go.uber.org/zap"

	"github.com/kucow/golang-grpc-base-project/internal/repo"
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

func (inter *Interceptor) AuthInterceptor() grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		inter.hasAuthorization()

		// token, err := grpc_auth.AuthFromMD(ctx, common.TokenType)
		// if err != nil {
		// 	return nil, err
		// }
		//
		// _, err = common.VerifyToken(inter.log, token)
		// if err != nil {
		// 	return nil, err
		// }
		//
		// // userID := cast.ToUint64(verifiedToken.StandardClaims.Subject)
		// // if token != "123456789" {
		// // 	return nil, status.Errorf(codes.PermissionDenied, "bad token")
		// // }

		return context.WithValue(ctx, "some_context_marker", "marker_exists"), nil
	}
}

func (inter *Interceptor) hasAuthorization() {
	permissions := inter.persist.FindAllPermissions()
	if len(permissions) == 0 {
		return
	}

	var authorize []map[string][]string
	for _, permission := range permissions {
		var authorizeRoles []string

		edgeRoles := permission.Edges.Roles
		if len(edgeRoles) > 0 {
			for _, role := range edgeRoles {
				authorizeRoles = append(authorizeRoles, role.Slug)
			}
		}

		authorize = append(authorize, map[string][]string{
			permission.Slug: authorizeRoles,
		})
	}

	inter.log.Info("Permission",
		zap.Any("authorize", authorize),
	)
}
