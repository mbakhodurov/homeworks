package v1

import (
	"context"
	"errors"

	"github.com/mbakhodurov/homeworks/week2/order/internal/model"
	order_v1 "github.com/mbakhodurov/homeworks/week2/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateOrderRes, error) {

	order, err := a.OrderService.Create(ctx, req.UserUUID, req.PartUuids)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrPartsNotFound):
			return &order_v1.BadRequestError{
				Code:    400,
				Message: "Не все необходимые детали найдены",
			}, nil
		default:
			return nil, err
		}
	}

	return &order_v1.CreateOrderResponse{
		OrderUUID:  order.Order_uuid,
		TotalPrice: float32(order.Total_price),
	}, nil
}
