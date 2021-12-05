package validator

import (
	"fmt"

	"github.com/AmreeshTyagi/goldflake"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kucow/golang-grpc-base-project/internal/repo"
	commonproto "github.com/kucow/golang-grpc-base-project/pkg/proto/v1/common"
)

const (
	FmtValidate = "%v: %v"
)

var (
	IsULID = validation.NewStringRule(ValidateUUID, "must be a valid id")
)

type Validator struct {
	log     *zap.Logger
	persist repo.Persist
}

func NewValidator(log *zap.Logger, persist repo.Persist) *Validator {
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

func ValidateUUID(id string) bool {
	decompose := goldflake.Decompose(cast.ToUint64(id))
	return len(decompose) > 0
}

func (val *Validator) ValidateCommonID(in *commonproto.UUIDRequest) error {
	err := validation.ValidateStruct(in,
		// Validate id
		validation.Field(&in.Id,
			validation.Required,
			validation.Length(5, 100),
			IsULID,
		),
	)

	return ValidateError(err)
}
