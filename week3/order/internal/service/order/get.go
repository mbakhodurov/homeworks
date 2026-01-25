package order

import (
	"context"

	"github.com/mbakhodurov/homeworks/week3/order/internal/model"
)

func (s *service) Get(ctx context.Context, order_uuid string) (model.Order, error) {
	res, err := s.orderRepo.Get(ctx, order_uuid)
	if err != nil {
		return model.Order{}, err
	}
	return res, nil
}
