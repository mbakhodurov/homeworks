package grpc

import (
	"context"

	"github.com/mbakhodurov/homeworks/week2/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID string, userUUID string, paymentMethod model.PaymentMethod) (string, error)
}
