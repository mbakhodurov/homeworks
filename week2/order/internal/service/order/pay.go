package order

import (
	"context"
	"time"

	"github.com/mbakhodurov/homeworks/week2/order/internal/model"
)

func (s *service) Pay(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) (string, error) {
	order, err := s.orderRepo.Get(ctx, orderUUID)
	if err != nil {
		return "", model.ErrOrderNotFound
	}

	if order.Status != model.StatusPendingPayment {
		switch order.Status {
		case model.StatusPaid:
			return "", model.ErrOrderAlreadyPaid
		case model.StatusCancelled:
			return "", model.ErrOrderCancelled
		default:
			return "", model.ErrInvalidOrderStatus
		}
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, orderUUID, order.User_uuid, paymentMethod)
	if err != nil {
		return "", err
	}

	now := time.Now()
	status := model.StatusPaid

	updateInfo := model.OrderUpdateInfo{
		Status:           &status,
		Transaction_uuid: &transactionUUID,
		Payment_method:   &paymentMethod,
		Updated_at:       &now,
	}

	if err := s.orderRepo.Update(ctx, orderUUID, updateInfo); err != nil {
		return "", err
	}

	return transactionUUID, nil
}
