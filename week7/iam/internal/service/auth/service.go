package auth

import (
	"time"

	"github.com/mbakhodurov/homeworks/week7/iam/internal/repository"
	iamService "github.com/mbakhodurov/homeworks/week7/iam/internal/service"
)

var _ iamService.AuthService = (*service)(nil)

type service struct {
	sessionRepository repository.SessionRepository
	userRepository    repository.UserRepository
	cacheTTL          time.Duration
}

func NewService(sessionRepository repository.SessionRepository, userRepository repository.UserRepository, cacheTTL time.Duration) *service {
	return &service{
		sessionRepository: sessionRepository,
		userRepository:    userRepository,
		cacheTTL:          cacheTTL,
	}
}
