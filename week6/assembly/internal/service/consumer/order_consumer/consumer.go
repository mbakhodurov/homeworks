package order_consumer

import (
	"context"

	kafkaConverter "github.com/mbakhodurov/homeworks/week6/assembly/internal/converter/kafka"
	def "github.com/mbakhodurov/homeworks/week6/assembly/internal/service"
	"go.uber.org/zap"

	"github.com/mbakhodurov/homeworks/week6/platform/pkg/kafka"
	"github.com/mbakhodurov/homeworks/week6/platform/pkg/logger"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	orderConsumer    kafka.Consumer
	orderPaidDecoder kafkaConverter.OrderPaidDecoder
	orderProducer    def.OrderProducerService
}

func NewService(orderConsumer kafka.Consumer, orderPaidDecoder kafkaConverter.OrderPaidDecoder, orderProducer def.OrderProducerService) *service {
	return &service{
		orderConsumer:    orderConsumer,
		orderPaidDecoder: orderPaidDecoder,
		orderProducer:    orderProducer,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order consumer service")

	err := s.orderConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Failed to start order consumer service", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Order consumer service successfully started")
	return nil
}
