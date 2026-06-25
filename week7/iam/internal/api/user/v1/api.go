package v1

import (
	"github.com/mbakhodurov/homeworks/week7/iam/internal/service"
	user_v1 "github.com/mbakhodurov/homeworks/week7/shared/pkg/proto/user/v1"
)

type api struct {
	user_v1.UnimplementedUserServiceServer

	userService service.UserService
}

func NewAPI(userService service.UserService) *api {
	return &api{
		userService: userService,
	}
}
