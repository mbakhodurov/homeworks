package order

import (
	"context"
	"time"

	"github.com/mbakhodurov/homeworks/week2/order/internal/model"
)

func (s *service) CancelOrder(ctx context.Context, orderUUID string) error {
	order, err := s.orderRepo.Get(ctx, orderUUID)
	if err != nil {
		return model.ErrOrderNotFound
	}

	if order.Status == model.StatusPaid {
		return model.ErrOrderAlreadyPaid
	}

	if order.Status == model.StatusCancelled {
		return model.ErrOrderCancelled
	}

	order.Status = model.StatusCancelled

	now := time.Now()
	status := model.StatusCancelled

	updateInfo := model.OrderUpdateInfo{
		Status:     &status,
		Deleted_at: &now,
	}

	if err := s.orderRepo.Update(ctx, orderUUID, updateInfo); err != nil {
		return err
	}

	return nil
}
