package validator_handler

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/wire"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xdorro/golang-grpc-base-project/internal/repo"
)

// ProviderSet is validator providers.
var ProviderSet = wire.NewSet(NewValidator)

type Validator struct {
	log  *zap.Logger
	repo *repo.Repo
}

func NewValidator(log *zap.Logger, repo *repo.Repo) ValidatorPersist {
	val := &Validator{
		log:  log,
		repo: repo,
	}

	return val
}

// ValidateError validate payload error if not nil
func (val *Validator) ValidateError(err error) error {
	if err != nil {
		if e, ok := err.(validation.Errors); ok {
			for name, value := range e {
				return status.New(codes.InvalidArgument, fmt.Sprintf("%v: %v", name, value)).Err()
			}
		}

		return status.New(codes.InvalidArgument, err.Error()).Err()
	}

	return nil
}
