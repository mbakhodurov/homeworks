package v1

import (
	"context"
	"errors"

	"github.com/mbakhodurov/homeworks/week5/order/internal/converter"
	"github.com/mbakhodurov/homeworks/week5/order/internal/model"
	order_v1 "github.com/mbakhodurov/homeworks/week5/shared/pkg/openapi/order/v1"
)

func (a *api) GetAllOrders(ctx context.Context) (order_v1.GetAllOrdersRes, error) {
	orders, err := a.orderService.GetAll(ctx)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
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

	orderDtoList := make([]order_v1.OrderDto, 0, len(orders))

	for _, v := range orders {
		dto := converter.OrderModelToOpenApiModelOrder(v)
		orderDtoList = append(orderDtoList, dto)
	}

	return &order_v1.GetAllOrderResponse{
		OrderDto: orderDtoList,
	}, nil
}
