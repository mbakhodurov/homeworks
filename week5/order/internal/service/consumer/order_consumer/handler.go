package orderconsumer

import (
	"context"
	"time"

	"github.com/mbakhodurov/homeworks/week5/order/internal/model"
	"github.com/mbakhodurov/homeworks/week5/platform/pkg/kafka"
	"github.com/mbakhodurov/homeworks/week5/platform/pkg/logger"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func (s *service) OrderHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderShipAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode order event", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("event_uuid", event.EventUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.Any("build_time_sec", event.BuildTimeSec),
	)

	order, err := s.orderRepository.GetOrderByUUID(ctx, event.OrderUUID)
	if err != nil {
		logger.Error(ctx, "Failed to get order", zap.String("order_uuid", event.OrderUUID), zap.Error(err))
		return err
	}

	updateInfo := model.OrderUpdateInfo{
		Transaction_uuid: order.TransactionUUID,
		Payment_method:   order.PaymentMethod,
		Status:           lo.ToPtr(model.StatusCompleted),
		Updated_at:       lo.ToPtr(time.Now()),
	}

	if err := s.orderRepository.Update(ctx, event.OrderUUID, updateInfo); err != nil {
		logger.Error(ctx, "Failed to update order", zap.String("order_uuid", event.OrderUUID), zap.Error(err))
		return err
	}
	return nil
}
