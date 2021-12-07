package validator

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xdorro/golang-grpc-base-project/internal/persist"
)

const (
	FmtValidate = "%v: %v"
)

type Validator struct {
	log     *zap.Logger
	persist persist.Persist
}

func NewValidator(log *zap.Logger, persist persist.Persist) *Validator {
	val := &Validator{
		log:     log,
		persist: persist,
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
