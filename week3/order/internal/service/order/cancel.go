package order

import (
	"context"
	"time"

	"github.com/mbakhodurov/homeworks/week3/order/internal/model"
	"github.com/samber/lo"
)

func (s *service) Cancel(ctx context.Context, order_uuid string) error {
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

	status := model.StatusCancelled
	updateInfo := model.OrderUpdateInfo{
		Status:     &status,
		Updated_at: lo.ToPtr(time.Now()),
	}

	if err := s.orderRepo.Update(ctx, order_uuid, updateInfo); err != nil {
		return err
	}

	return nil

}
