package validator

import (
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/wire"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xdorro/golang-grpc-base-project/internal/repo"
)

// ProviderSet is validator providers.
var ProviderSet = wire.NewSet(NewValidator)

const (
	FmtValidate = "%v: %v"
)

var (
	PhoneNumber = regexp.MustCompile(`^(([03+[2-9]|05+[6|8|9]|07+[0|6|7|8|9]|08+[1-9]|09+[1-4|6-9]]){3})+[0-9]{7}\b$`)
)

type Validator struct {
	log  *zap.Logger
	repo *repo.Repo
}

func NewValidator(log *zap.Logger, repo *repo.Repo) *Validator {
	val := &Validator{
		log:  log,
		repo: repo,
	}

	return val
}

// ValidateError validate payload error if not nil
func ValidateError(err error) error {
	if err != nil {
		if e, ok := err.(validation.Errors); ok {
			for name, value := range e {
				return status.New(codes.InvalidArgument, fmt.Sprintf(FmtValidate, name, value)).Err()
			}
		}

		return status.New(codes.InvalidArgument, err.Error()).Err()
	}

	return nil
}
