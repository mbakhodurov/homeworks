package orderpaidconsumer

import (
	"context"

	kafkaConverter "github.com/mbakhodurov/homeworks/week7/notification/internal/converter/kafka"
	def "github.com/mbakhodurov/homeworks/week7/notification/internal/service"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/kafka"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/logger"
	"go.uber.org/zap"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	orderConsumer    kafka.Consumer
	orderPaidDecoder kafkaConverter.OrderPaidDecoder
	telegramService  def.TelegramService
}

func NewService(orderConsumer kafka.Consumer, orderPaidDecoder kafkaConverter.OrderPaidDecoder, telegramService def.TelegramService) *service {
	return &service{
		orderConsumer:    orderConsumer,
		orderPaidDecoder: orderPaidDecoder,
		telegramService:  telegramService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order consumer service")

	err := s.orderConsumer.Consume(ctx, s.OrderPaidHandler)
	if err != nil {
		logger.Error(ctx, "Failed to start order paid consumer service", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Order paid consumer service successfully started")
	return nil
}
