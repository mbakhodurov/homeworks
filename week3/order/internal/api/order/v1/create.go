package v1

import (
	"context"
	"errors"

	"github.com/mbakhodurov/homeworks/week3/order/internal/model"
	order_v1 "github.com/mbakhodurov/homeworks/week3/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateOrderRes, error) {
	res, err := a.orderService.Create(ctx, req.UserUUID, req.PartUuids)
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
		OrderUUID:  res.OrderUUID,
		TotalPrice: float32(res.TotalPrice),
	}, nil
}
