package handler

import (
	"github.com/google/wire"

	auth_handler "github.com/xdorro/golang-grpc-base-project/internal/handler/auth"
	validator_handler "github.com/xdorro/golang-grpc-base-project/internal/handler/validator"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	wire.Struct(new(Handler), "*"),
	validator_handler.ProviderSet,
	auth_handler.ProviderSet,
)

type Handler struct {
	validator_handler.ValidatorPersist
	auth_handler.AuthPersist
}
