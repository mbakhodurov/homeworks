package order_producer

import (
	"context"

	"github.com/mbakhodurov/homeworks/week7/order/internal/model"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/logger"
	events_v1 "github.com/mbakhodurov/homeworks/week7/shared/pkg/proto/events/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (s *service) ProduceOrderPaid(ctx context.Context, event model.OrderPaid) error {
	msg := &events_v1.OrderPaid{
		EventUuid:       event.EventUUID,
		OrderUuid:       event.OrderUUID,
		UserUuid:        event.UserUUID,
		PaymentMethod:   event.PaymentMethod,
		TransactionUuid: event.TransactionUUID,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal order paid")
		return err
	}

	err = s.orderProducer.Send(ctx, []byte(event.OrderUUID), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish order paid", zap.Any("event", event), zap.Error(err))
		return err
	}

	return nil
}
