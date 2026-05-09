package user

import (
	"context"

	"github.com/mbakhodurov/homeworks/week6/iam/internal/model"
)

func (s *service) GetAll(ctx context.Context) ([]model.User, error) {
	user, err := s.userRepository.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}
