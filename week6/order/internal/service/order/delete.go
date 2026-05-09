package order

import (
	"context"

	"github.com/mbakhodurov/homeworks/week6/order/internal/model"
)

func (s *service) Delete(ctx context.Context, order_uuid string) error {
	res, err := s.orderRepo.GetOrderByUUID(ctx, order_uuid)
	if err != nil {
		return err
	}

	if res.Status == model.StatusPaid {
		return model.ErrOrderAlreadyPaid
	}

	if res.Status == model.StatusCancelled {
		return model.ErrOrderCancelled
	}

	err = s.orderRepo.Delete(ctx, order_uuid)
	if err != nil {
		return err
	}
	return nil
}
