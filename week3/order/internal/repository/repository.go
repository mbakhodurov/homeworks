package repository

import (
	"context"

	"github.com/mbakhodurov/homeworks/week3/order/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order model.Order) (int64, error)
	Get(ctx context.Context, uuid string) (model.Order, error)
	GetAll(ctx context.Context) ([]model.Order, error)
	Update(ctx context.Context, uuid string, newData model.OrderUpdateInfo) error
	Delete(ctx context.Context, uuid string) error
}
