package v1

import (
	"context"
	"errors"

	"github.com/mbakhodurov/homeworks/week3/order/internal/converter"
	"github.com/mbakhodurov/homeworks/week3/order/internal/model"
	order_v1 "github.com/mbakhodurov/homeworks/week3/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrderByUUID(ctx context.Context, params order_v1.GetOrderByUUIDParams) (order_v1.GetOrderByUUIDRes, error) {
	order, err := a.orderService.Get(ctx, params.OrderUUID)
	// fmt.Println("order:", order)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &order_v1.NotFoundError{
				Code:    404,
				Message: "Заказов пока что нету",
			}, nil
		}
		return &order_v1.InternalServerError{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	return &order_v1.GetOrderResponse{
		OrderDto: converter.OrderModelToOrderDTO(order),
	}, nil
}
