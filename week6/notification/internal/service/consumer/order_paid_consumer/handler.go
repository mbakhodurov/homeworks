package orderpaidconsumer

import (
	"context"

	"github.com/mbakhodurov/homeworks/week6/platform/pkg/kafka"
	"github.com/mbakhodurov/homeworks/week6/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) OrderPaidHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode order paid event", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("event_uuid", event.EventUUID),
		zap.String("payment_method", event.PaymentMethod),
		zap.String("transaction_uuid", event.TransactionUUID),
	)

	if err := s.telegramService.SendOrderPaidNotification(ctx, event); err != nil {
		logger.Error(ctx, "Failed to send order paid notification", zap.Error(err))
		return err
	}
	return nil
}
