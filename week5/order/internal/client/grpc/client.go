package client

import (
	"context"

	"github.com/mbakhodurov/homeworks/week5/order/internal/model"
)

type InventoryClient interface {
	ListPartInventory(ctx context.Context, filter model.InventoryFilter) ([]model.Inventory, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, paymentMethod model.PaymentMethod, orderUUID, userUUID string) (string, error)
}
