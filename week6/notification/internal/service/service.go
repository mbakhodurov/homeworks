package service

import (
	"context"

	"github.com/mbakhodurov/homeworks/week6/notification/internal/model"
)

type TelegramService interface {
	SendShipAssembledNotification(ctx context.Context, event model.ShipAssembled) error
	SendOrderPaidNotification(ctx context.Context, event model.OrderPaid) error
}

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}
