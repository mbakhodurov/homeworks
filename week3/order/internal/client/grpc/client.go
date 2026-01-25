package grpc

import (
	"context"

	"github.com/mbakhodurov/homeworks/week3/order/internal/model"
)

type InventoryClient interface {
	InventoryPart(ctx context.Context, filter model.InventoryFilter) ([]model.Inventory, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, paymentMethod model.PaymentMethod, orderUUID, userUUID string) (string, error)
}
