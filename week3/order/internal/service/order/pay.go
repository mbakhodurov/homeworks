package order

import (
	"context"
	"time"

	"github.com/mbakhodurov/homeworks/week3/order/internal/model"
	"github.com/samber/lo"
)

const paymentTimeout = 3 * time.Second

func (s *service) Pay(ctx context.Context, order_uuid string, payment_method model.PaymentMethod) (string, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, paymentTimeout)
	defer cancel()

	order, err := s.orderRepo.Get(ctx, order_uuid)
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

	transactionUUID, err := s.paymentClient.PayOrder(ctxWithTimeout, payment_method, order_uuid, order.UserUUID)
	if err != nil {
		return "", err
	}
	status := model.StatusPaid
	updateInfo := model.OrderUpdateInfo{
		Status:           &status,
		Transaction_uuid: &transactionUUID,
		Payment_method:   &payment_method,
		Updated_at:       lo.ToPtr(time.Now()),
	}

	if err := s.orderRepo.Update(ctx, order_uuid, updateInfo); err != nil {
		return "", err
	}

	return transactionUUID, nil
}
