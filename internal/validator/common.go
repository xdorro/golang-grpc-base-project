package validator

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	common_proto "github.com/xdorro/golang-grpc-base-project/api/proto/v1/common"
)

func (val *Validator) ValidateCommonID(in *common_proto.UUIDRequest) error {
	err := validation.ValidateStruct(in,
		// Validate id
		validation.Field(&in.Id,
			validation.Required,
			is.Int,
		),
	)

	return ValidateError(err)
}
