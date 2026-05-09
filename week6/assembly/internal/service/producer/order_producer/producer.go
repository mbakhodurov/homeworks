package order_producer

import (
	"context"

	"github.com/mbakhodurov/homeworks/week6/assembly/internal/model"
	"github.com/mbakhodurov/homeworks/week6/platform/pkg/logger"
	events_v1 "github.com/mbakhodurov/homeworks/week6/shared/pkg/proto/events/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (s *service) ProduceShipAssembled(ctx context.Context, event model.ShipAssembled) error {
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

	if err := s.orderProducer.Send(ctx, []byte(event.EventUUID), payload); err != nil {
		logger.Error(ctx, "failed to publish order assembled", zap.Error(err))
		return err
	}
	return nil
}
