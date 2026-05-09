package v1

import (
	"github.com/mbakhodurov/homeworks/week6/iam/internal/service"
	auth_v1 "github.com/mbakhodurov/homeworks/week6/shared/pkg/proto/auth/v1"
)

type api struct {
	auth_v1.UnimplementedAuthServiceServer

	authService service.AuthService
}

func NewAPI(authService service.AuthService) *api {
	return &api{
		authService: authService,
	}
}
