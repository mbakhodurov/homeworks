package v1

import (
	"github.com/mbakhodurov/homeworks/week5/order/internal/service"
	order_v1 "github.com/mbakhodurov/homeworks/week5/shared/pkg/openapi/order/v1"
)

var _ order_v1.Handler = (*api)(nil)

type api struct {
	orderService service.OrderService
}

func NewApi(orderService service.OrderService) *api {
	return &api{
		orderService: orderService,
	}
}
