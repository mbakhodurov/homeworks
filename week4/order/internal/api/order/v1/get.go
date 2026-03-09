package v1

import (
	"context"
	"errors"

	"github.com/mbakhodurov/homeworks/week4/order/internal/converter"
	"github.com/mbakhodurov/homeworks/week4/order/internal/model"
	"github.com/mbakhodurov/homeworks/week4/platform/pkg/logger"
	order_v1 "github.com/mbakhodurov/homeworks/week4/shared/pkg/openapi/order/v1"
	"go.uber.org/zap"
)

func (a *api) GetOrderByUUID(ctx context.Context, params order_v1.GetOrderByUUIDParams) (order_v1.GetOrderByUUIDRes, error) {
	res, err := a.orderService.Get(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &order_v1.NotFoundError{
				Code:    404,
				Message: "Заказов пока что нету",
			}, nil
		}
		logger.Error(ctx, "ошибка", zap.Error(err))
		return &order_v1.InternalServerError{
			Code:    500,
			Message: err.Error(),
		}, nil
	}
	return &order_v1.GetOrderResponse{
		OrderDto: converter.OrderModelToOrderDTO(res),
	}, nil
}
