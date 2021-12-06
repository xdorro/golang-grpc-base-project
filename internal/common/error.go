package common

import (
	"reflect"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	Success = status.New(codes.OK, "success")

	UserNotExist      = status.New(codes.InvalidArgument, "user doesn't exist")
	EmailAlreadyExist = status.New(codes.InvalidArgument, "user email already exist")
	EmailNotExist     = status.New(codes.InvalidArgument, "user email doesn't exist")

	TokenInvalid      = status.New(codes.Unauthenticated, "token is invalid")
	PasswordIncorrect = status.New(codes.InvalidArgument, "password is incorrect")

	RoleAlreadyExist = status.New(codes.InvalidArgument, "role slug already exist")
	RoleNotExist     = status.New(codes.InvalidArgument, "role doesn't exist")

	PermissionAlreadyExist = status.New(codes.InvalidArgument, "permission slug already exist")
	PermissionNotExist     = status.New(codes.InvalidArgument, "permission doesn't exist")
)

func CheckError(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}

func CompareError(err1, err2 error) bool {
	return reflect.DeepEqual(CheckError(err1), CheckError(err2))
}
