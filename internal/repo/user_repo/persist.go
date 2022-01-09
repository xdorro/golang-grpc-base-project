package user_repo

import (
	"github.com/xdorro/golang-grpc-base-project/api/ent"
)

type UserPersist interface {
	FindAllUsers() []*ent.User
	FindUserByID(id uint64) (*ent.User, error)
	FindUserByEmail(email string) (*ent.User, error)
	ExistUserByID(id uint64) bool
	ExistUserByEmail(email string) bool
	CreateUser(user *ent.User) error
	UpdateUser(user *ent.User) error
	DeleteUser(id uint64) error
	SoftDeleteUser(id uint64) error
}
