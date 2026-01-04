package order

import (
	"context"

	"github.com/mbakhodurov/homeworks/week2/order/internal/model"
)

func (s *service) Get(ctx context.Context, OrderUUID string) (model.Order, error) {
	order, err := s.orderRepo.Get(ctx, OrderUUID)
	if err != nil {
		return model.Order{}, err
	}
	return order, nil
}
