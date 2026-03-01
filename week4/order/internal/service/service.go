package service

import (
	"context"

	"github.com/mbakhodurov/homeworks/week4/order/internal/model"
)

type OrderService interface {
	Create(ctx context.Context, user_uuid string, partUUIDs []string) (*model.Order, error)
}
