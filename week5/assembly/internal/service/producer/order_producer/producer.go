package orderproducer

import (
	"context"

	"github.com/mbakhodurov/homeworks/week5/assembly/internal/model"
	"github.com/mbakhodurov/homeworks/week5/platform/pkg/kafka"
	"github.com/mbakhodurov/homeworks/week5/platform/pkg/logger"
	events_v1 "github.com/mbakhodurov/homeworks/week5/shared/pkg/proto/events/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	def "github.com/mbakhodurov/homeworks/week5/assembly/internal/service"
)

var _ def.OrderProducerService = (*service)(nil)

type service struct {
	orderProducer kafka.Producer
}

func NewService(orderProducer kafka.Producer) *service {
	return &service{
		orderProducer: orderProducer,
	}
}

func (p *service) ProduceShipAssembled(ctx context.Context, event model.ShipAssembled) error {
	msg := &events_v1.ShipAssembled{
		EventUuid:    event.EventUUID,
		OrderUuid:    event.OrderUUID,
		UserUuid:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal OrderAssembled", zap.Error(err))
		return err
	}

	if err := p.orderProducer.Send(ctx, []byte(event.EventUUID), payload); err != nil {
		logger.Error(ctx, "failed to publish order assembled", zap.Error(err))
		return err
	}
	return nil
}
