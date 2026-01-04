package order

import (
	"context"

	"github.com/mbakhodurov/homeworks/week2/order/internal/model"
)

func (s *service) GetAll(ctx context.Context) ([]model.Order, error) {
	orders, err := s.orderRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
