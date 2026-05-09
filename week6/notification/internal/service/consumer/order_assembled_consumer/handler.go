package orderassembledconsumer

import (
	"context"

	"github.com/mbakhodurov/homeworks/week6/platform/pkg/kafka"
	"github.com/mbakhodurov/homeworks/week6/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) OrderAssembledHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.shipAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode order assembled event", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.Int64("build_time_sec", event.BuildTimeSec),
	)

	err = s.telegramService.SendShipAssembledNotification(ctx, event)
	if err != nil {
		logger.Error(ctx, "Failed to send ship assembled notification", zap.Error(err))
		return err
	}

	return nil
}
