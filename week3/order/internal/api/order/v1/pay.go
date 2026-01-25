package v1

import (
	"context"
	"errors"

	"github.com/mbakhodurov/homeworks/week3/order/internal/converter"
	"github.com/mbakhodurov/homeworks/week3/order/internal/model"
	order_v1 "github.com/mbakhodurov/homeworks/week3/shared/pkg/openapi/order/v1"
)

func (a *api) PaymentOrder(ctx context.Context, req *order_v1.PayOrderRequest, params order_v1.PaymentOrderParams) (order_v1.PaymentOrderRes, error) {
	transactionUUID, err := a.orderService.Pay(ctx, params.OrderUUID, converter.PaymentMethodProtoToPaymentModel(req.PaymentMethod))
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return &order_v1.NotFoundError{
				Code:    404,
				Message: "Заказ не найден",
			}, nil
		case errors.Is(err, model.ErrOrderAlreadyPaid):
			return &order_v1.ConflictError{
				Code:    409,
				Message: "Заказ уже оплачен",
			}, nil
		case errors.Is(err, model.ErrOrderCancelled):
			return &order_v1.ConflictError{
				Code:    409,
				Message: "Нельзя оплатить отмененный заказ",
			}, nil
		default:
			return nil, err
		}
	}

	return &order_v1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}
