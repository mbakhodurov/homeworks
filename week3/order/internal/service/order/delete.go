package order

import (
	"context"

	"github.com/mbakhodurov/homeworks/week3/order/internal/model"
)

func (s *service) Delete(ctx context.Context, order_uuid string) error {
	order, err := s.orderRepo.Get(ctx, order_uuid)
	if err != nil {
		return err
	}

	if order.Status == model.StatusPaid {
		return model.ErrOrderAlreadyPaid
	}

	if order.Status == model.StatusCancelled {
		return model.ErrOrderCancelled
	}

	if err := s.orderRepo.Delete(ctx, order_uuid); err != nil {
		return err
	}

	return nil

}
