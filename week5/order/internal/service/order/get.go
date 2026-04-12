package order

import (
	"context"

	"github.com/mbakhodurov/homeworks/week5/order/internal/model"
)

func (s *service) Get(ctx context.Context, orderUUID string) (model.Order, error) {
	res, err := s.orderRepo.GetOrderByUUID(ctx, orderUUID)
	if err != nil {
		return model.Order{}, err
	}
	return res, nil
}
