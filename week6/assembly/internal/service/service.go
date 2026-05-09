package service

import (
	"context"

	"github.com/mbakhodurov/homeworks/week6/assembly/internal/model"
)

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type OrderProducerService interface {
	ProduceShipAssembled(ctx context.Context, event model.ShipAssembled) error
}
