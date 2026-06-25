package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mbakhodurov/homeworks/week7/order/internal/converter"
	"github.com/mbakhodurov/homeworks/week7/order/internal/model"
	"github.com/samber/lo"
)

const paymentTimeout = 3 * time.Second

func (s *service) Pay(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) (string, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, paymentTimeout)
	defer cancel()

	order, err := s.orderRepo.GetOrderByUUID(ctx, orderUUID)
	if err != nil {
		return "", model.ErrOrderNotFound
	}

	if order.Status != model.StatusPendingPayment {
		switch order.Status {
		case model.StatusPaid:
			return "", model.ErrOrderAlreadyPaid
		case model.StatusCancelled:
			return "", model.ErrOrderCancelled
		case model.StatusCompleted:
			return "", model.ErrOrderCompleted
		default:
			fmt.Println("fdsa", order.Status)
			return "", model.ErrInvalidOrderStatus
		}
	}

	transactionUUID, err := s.paymenClient.PayOrder(ctxWithTimeout, order.OrderUUID, order.UserUUID, paymentMethod)
	if err != nil {
		return "", err
	}

	status := model.StatusPaid
	modelUpdated := model.OrderUpdateInfo{
		Status:           &status,
		Transaction_uuid: &transactionUUID,
		Payment_method:   &paymentMethod,
		Updated_at:       lo.ToPtr(time.Now()),
	}

	if err := s.orderRepo.Update(ctx, orderUUID, modelUpdated); err != nil {
		return "", err
	}

	err = s.orderProducerService.ProduceOrderPaid(ctx, model.OrderPaid{
		EventUUID:       uuid.NewString(),
		OrderUUID:       orderUUID,
		UserUUID:        order.UserUUID,
		PaymentMethod:   converter.PaymentMethodDtoToString(paymentMethod),
		TransactionUUID: transactionUUID,
	})
	if err != nil {
		return "", err
	}
	return transactionUUID, nil
}
