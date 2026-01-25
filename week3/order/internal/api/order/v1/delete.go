package v1

import (
	"context"
	"errors"

	"github.com/mbakhodurov/homeworks/week3/order/internal/model"
	order_v1 "github.com/mbakhodurov/homeworks/week3/shared/pkg/openapi/order/v1"
)

func (a *api) DeleteOrder(ctx context.Context, params order_v1.DeleteOrderParams) (order_v1.DeleteOrderRes, error) {
	if err := a.orderService.Delete(ctx, params.OrderUUID); err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return &order_v1.NotFoundError{
				Code:    404,
				Message: "Заказ не найден",
			}, nil
		case errors.Is(err, model.ErrOrderAlreadyPaid):
			return &order_v1.ConflictError{
				Code:    409,
				Message: "Оплаченный заказ нельзя удалить",
			}, nil
		case errors.Is(err, model.ErrOrderCancelled):
			return &order_v1.ConflictError{
				Code:    409,
				Message: "Отмененный заказ нельзя удалить",
			}, nil
		default:
			return nil, err

		}
	}

	return &order_v1.DeleteOrderResponse{
		Properties: order_v1.NewOptString("Успешно удален заказ"),
	}, nil
}
