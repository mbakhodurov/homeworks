package repository

import (
	"context"

	"github.com/mbakhodurov/homeworks/week2/order/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order model.Order) error
	GetAll(ctx context.Context) ([]model.Order, error)
	Get(ctx context.Context, OrderUUID string) (model.Order, error)
	Update(ctx context.Context, uuid string, info model.OrderUpdateInfo) error
}
