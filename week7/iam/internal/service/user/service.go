package user

import (
	"github.com/mbakhodurov/homeworks/week7/iam/internal/repository"
	iamService "github.com/mbakhodurov/homeworks/week7/iam/internal/service"
)

var _ iamService.UserService = (*service)(nil)

type service struct {
	userRepository repository.UserRepository
}

func NewService(userRepository repository.UserRepository) *service {
	return &service{userRepository: userRepository}
}
