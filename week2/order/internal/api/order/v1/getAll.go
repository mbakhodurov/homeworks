package v1

import (
	"context"
	"errors"

	"github.com/mbakhodurov/homeworks/week2/order/internal/converter"
	"github.com/mbakhodurov/homeworks/week2/order/pkg/models"
	order_v1 "github.com/mbakhodurov/homeworks/week2/shared/pkg/openapi/order/v1"
)

func (a *api) GetAllOrders(ctx context.Context) (order_v1.GetAllOrdersRes, error) {
	orders, err := a.OrderService.GetAll(ctx)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
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

	for _, order := range orders {
		dto := converter.OrderToDTO(order)
		orderDtoList = append(orderDtoList, dto)
	}
	return &order_v1.GetAllOrderResponse{
		OrderDto:   orderDtoList,
		TotalCount: float64(len(orderDtoList)),
	}, nil
}
