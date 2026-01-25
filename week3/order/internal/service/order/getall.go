package order

import (
	"context"

	"github.com/mbakhodurov/homeworks/week3/order/internal/model"
)

func (s *service) GetAll(ctx context.Context) ([]model.Order, error) {
	res, err := s.orderRepo.GetAll(ctx)
	if err != nil {
		return []model.Order{}, err
	}

	return res, nil
}
