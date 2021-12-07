package servicecommon

import (
	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/ent"
	"github.com/xdorro/golang-grpc-base-project/internal/persist"
	"github.com/xdorro/golang-grpc-base-project/pkg/validator"
)

type Service struct {
	Log       *zap.Logger
	Client    *ent.Client
	Persist   persist.Persist
	Validator *validator.Validator
}

func NewService(log *zap.Logger, client *ent.Client, persist persist.Persist, validator *validator.Validator) *Service {
	return &Service{Log: log, Client: client, Persist: persist, Validator: validator}
}
