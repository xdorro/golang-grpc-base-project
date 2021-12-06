package validator

import (
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kucow/golang-grpc-base-project/internal/common"
)

func TestValidateError(t *testing.T) {
	type Customer struct {
		Email string
	}

	c := Customer{
		Email: "Qiang Xue",
	}

	c2 := Customer{
		Email: "admin@example.com",
	}

	err := validation.ValidateStruct(c,
		// Validate email
		validation.Field(&c.Email,
			validation.Required,
			is.Email,
			validation.Length(5, 0),
		),
	)

	err2 := validation.ValidateStruct(&c,
		// Validate email
		validation.Field(&c.Email,
			validation.Required,
			is.Email,
			validation.Length(5, 0),
		),
	)

	err3 := validation.ValidateStruct(&c2,
		// Validate email
		validation.Field(&c2.Email,
			validation.Required,
			is.Email,
			validation.Length(5, 0),
		),
	)

	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "empty args",
			args: args{},
		},
		{
			name: "error pointer",
			args: args{
				err: err,
			},
			wantErr: status.New(codes.InvalidArgument, "only a pointer to a struct can be validated").Err(),
		},
		{
			name: "error email",
			args: args{
				err: err2,
			},
			wantErr: status.New(codes.InvalidArgument, "Email: must be a valid email address").Err(),
		},
		{
			name: "success",
			args: args{
				err: err3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err = ValidateError(tt.args.err); !common.CompareError(err, tt.wantErr) {
				t.Errorf("ValidateError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
