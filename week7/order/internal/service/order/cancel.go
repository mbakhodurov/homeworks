package order

import (
	"context"
	"time"

	"github.com/mbakhodurov/homeworks/week7/order/internal/model"
	"github.com/samber/lo"
)

func (s *service) Cancel(ctx context.Context, orderUUID string) error {
	order, err := s.orderRepo.GetOrderByUUID(ctx, orderUUID)
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

	if err := s.orderRepo.Update(ctx, orderUUID, updateInfo); err != nil {
		return err
	}
	return nil
}
