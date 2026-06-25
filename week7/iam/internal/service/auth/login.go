package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mbakhodurov/homeworks/week7/iam/internal/model"
)

func (s *service) Login(ctx context.Context, login, password string) (sessionUUID string, err error) {
	user, err := s.userRepository.GetUserByLogin(ctx, login, password)
	if err != nil {
		fmt.Println(err.Error())
		return "", model.ErrInvalidCredentials
	}

	session := model.Session{
		UUID:      uuid.NewString(),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(s.cacheTTL),
	}

	err = s.sessionRepository.Create(ctx, session, user, s.cacheTTL)
	if err != nil {
		return "", err
	}

	return session.UUID, nil
}
