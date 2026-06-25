package client

import (
	"context"

	"github.com/mbakhodurov/homeworks/week7/order/internal/model"
)

type InventoryClient interface {
	ListPartInventory(ctx context.Context, filter model.InventoryFilter) ([]model.Inventory, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod model.PaymentMethod) (string, error)
}
