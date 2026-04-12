package order_assembled_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/mbakhodurov/homeworks/week5/notification/internal/converter/kafka"
	def "github.com/mbakhodurov/homeworks/week5/notification/internal/service"
	"github.com/mbakhodurov/homeworks/week5/platform/pkg/kafka"
	"github.com/mbakhodurov/homeworks/week5/platform/pkg/logger"
)

type service struct {
	orderConsumer        kafka.Consumer
	shipAssembledDecoder kafkaConverter.OrderAssembledDecoder
	telegramService      def.TelegramService
}

func NewService(orderConsumer kafka.Consumer, shipAssembledDecoder kafkaConverter.OrderAssembledDecoder, telegramService def.TelegramService) *service {
	return &service{
		orderConsumer:        orderConsumer,
		shipAssembledDecoder: shipAssembledDecoder,
		telegramService:      telegramService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting ship assembled consumer service")

	err := s.orderConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Failed to start ship assembled consumer service", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Ship assembly consumer service successfully started")
	return nil
}
