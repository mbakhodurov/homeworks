package service

import (
	"context"

	"github.com/mbakhodurov/homeworks/week2/order/internal/model"
)

type OrderService interface {
	Create(ctx context.Context, userUUID string, partsUUID []string) (*model.Order, error)
	GetAll(ctx context.Context) ([]model.Order, error)
	Get(ctx context.Context, OrderUUID string) (model.Order, error)
	CancelOrder(ctx context.Context, orderUUID string) error
	Pay(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) (string, error)
}
