package repository

import (
	"context"

	"github.com/mbakhodurov/homeworks/week7/order/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order model.Order) (int64, error)
	GetOrderByUUID(ctx context.Context, orderUUID string) (model.Order, error)
	GetAll(ctx context.Context) ([]model.Order, error)
	Delete(ctx context.Context, order_uuid string) error
	Update(ctx context.Context, order_uuid string, updateInfo model.OrderUpdateInfo) error
}
