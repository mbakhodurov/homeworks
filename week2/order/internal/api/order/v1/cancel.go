package v1

import (
	"context"
	"errors"

	"github.com/mbakhodurov/homeworks/week2/order/internal/model"
	order_v1 "github.com/mbakhodurov/homeworks/week2/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params order_v1.CancelOrderParams) (order_v1.CancelOrderRes, error) {
	err := a.OrderService.CancelOrder(ctx, params.OrderUUID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return &order_v1.NotFoundError{
				Code:    404,
				Message: "Not found",
			}, nil
		case errors.Is(err, model.ErrOrderAlreadyPaid):
			return &order_v1.ConflictError{
				Code:    409,
				Message: "Заказ уже оплачен и не может быть отменён",
			}, nil

		case errors.Is(err, model.ErrOrderCancelled):
			return &order_v1.ConflictError{
				Code:    409,
				Message: "Заказ уже отменён",
			}, nil
		default:
			return nil, err
		}
	}
	return &order_v1.CancelOrderResponse{
		Properties: order_v1.NewOptString("Успешно отменен заказ"),
	}, nil
}
