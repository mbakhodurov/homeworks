package v1

import (
	"github.com/mbakhodurov/homeworks/week2/order/internal/service"
	order_v1 "github.com/mbakhodurov/homeworks/week2/shared/pkg/openapi/order/v1"
)

var _ order_v1.Handler = (*api)(nil)

type api struct {
	OrderService service.OrderService
}

func NewApi(OrderService service.OrderService) *api {
	return &api{
		OrderService: OrderService,
	}
}
