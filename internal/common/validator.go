package common

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commonproto "github.com/kucow/golang-grpc-base/pkg/proto/v1/common"
)

const (
	FmtValidate = "%v: %v"
)

var (
	IsULID = validation.NewStringRule(ValidateULID, "must be a valid id")
)

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

func ValidateULID(id string) bool {
	_, err := ulid.ParseStrict(id)
	return err == nil
}

func ValidateCommonID(in *commonproto.UUIDRequest) error {
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
