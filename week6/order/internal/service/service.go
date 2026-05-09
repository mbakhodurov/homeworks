package service

import (
	"context"

	"github.com/mbakhodurov/homeworks/week6/order/internal/model"
)

type OrderService interface {
	Create(ctx context.Context, userUUID string, partUUIDs []string) (*model.Order, error)
	Delete(ctx context.Context, order_uuid string) error
	Get(ctx context.Context, orderUUID string) (model.Order, error)
	GetAll(ctx context.Context) ([]model.Order, error)
	Pay(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) (string, error)
	Cancel(ctx context.Context, orderUUID string) error
}

type OrderProducerService interface {
	ProduceOrderPaid(ctx context.Context, event model.OrderPaid) error
}

type OrderConsumerService interface {
	RunConsumer(ctx context.Context) error
}
