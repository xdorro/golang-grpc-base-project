package validator_handler

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/xdorro/golang-grpc-base-project/api/proto/common"
)

func (val *Validator) ValidateCommonID(in *common_proto.UUIDRequest) error {
	err := validation.ValidateStruct(in,
		// Validate id
		validation.Field(&in.Id,
			validation.Required,
			is.Int,
		),
	)

	return val.ValidateError(err)
}

func (val *Validator) ValidateCommonSlug(in *common_proto.SlugRequest) error {
	err := validation.ValidateStruct(in,
		// Validate slug
		validation.Field(&in.Slug,
			validation.Required,
			// is.Int,
		),
	)

	return val.ValidateError(err)
}
