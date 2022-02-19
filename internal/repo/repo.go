package repo

import (
	"github.com/google/wire"

	"github.com/xdorro/golang-grpc-base-project/api/ent"
	permission_repo "github.com/xdorro/golang-grpc-base-project/internal/repo/permission"
	role_repo "github.com/xdorro/golang-grpc-base-project/internal/repo/role"
	user_repo "github.com/xdorro/golang-grpc-base-project/internal/repo/user"
)

// ProviderSet is repo providers.
var ProviderSet = wire.NewSet(
	wire.Struct(new(Repo), "*"),
	user_repo.ProviderSet,
	role_repo.ProviderSet,
	permission_repo.ProviderSet,
)

// Repo is a wrapper around an ent.Client that provides a convenient API for
type Repo struct {
	*ent.Client

	user_repo.UserPersist
	role_repo.RolePersist
	permission_repo.PermissionPersist
}

func (repo *Repo) Close() error {
	if repo.Client != nil {
		return repo.Client.Close()
	}

	return nil
}
