package repository

import (
	"context"

	"github.com/mbakhodurov/homeworks/week4/order/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order model.Order) (int64, error)
}
