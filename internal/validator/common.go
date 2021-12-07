package validator

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/xdorro/golang-grpc-base-project/proto/v1/common"
)

func (val *Validator) ValidateCommonID(in *commonproto.UUIDRequest) error {
	err := validation.ValidateStruct(in,
		// Validate id
		validation.Field(&in.Id,
			validation.Required,
			is.Int,
		),
	)

	return ValidateError(err)
}
