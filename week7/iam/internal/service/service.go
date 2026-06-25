package service

import (
	"context"

	"github.com/mbakhodurov/homeworks/week7/iam/internal/model"
)

type UserService interface {
	GetUser(ctx context.Context, userUUID string) (model.User, error)
	Register(ctx context.Context, userInfo model.UserInfo, password string) (userUUID string, err error)
	GetAll(ctx context.Context) ([]model.User, error)
}

type AuthService interface {
	Login(ctx context.Context, login, password string) (sessionUUID string, err error)
	Whoami(ctx context.Context, sessionUUID string) (model.Session, model.User, error)
}
