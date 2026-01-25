package v1

import (
	"github.com/mbakhodurov/homeworks/week3/order/internal/service"
	order_v1 "github.com/mbakhodurov/homeworks/week3/shared/pkg/openapi/order/v1"
)

type api struct {
	orderService service.OrderService
}

var _ order_v1.Handler = (*api)(nil)

func NewApi(orderService service.OrderService) *api {
	return &api{
		orderService: orderService,
	}
}
