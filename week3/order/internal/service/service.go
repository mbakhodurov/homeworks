package service

import (
	"context"

	"github.com/mbakhodurov/homeworks/week3/order/internal/model"
)

type OrderService interface {
	Create(ctx context.Context, user_uuid string, partUUIDs []string) (*model.Order, error)
	GetAll(ctx context.Context) ([]model.Order, error)
	Get(ctx context.Context, order_uuid string) (model.Order, error)
	Pay(ctx context.Context, order_uuid string, payment_method model.PaymentMethod) (string, error)
	Cancel(ctx context.Context, order_uuid string) error
	Delete(ctx context.Context, order_uuid string) error
}
