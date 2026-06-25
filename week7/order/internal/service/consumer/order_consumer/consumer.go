package orderconsumer

import (
	"context"

	orderDecoderInterface "github.com/mbakhodurov/homeworks/week7/order/internal/converter/kafka"
	"github.com/mbakhodurov/homeworks/week7/order/internal/repository"
	def "github.com/mbakhodurov/homeworks/week7/order/internal/service"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/kafka"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/logger"
	"go.uber.org/zap"
)

var _ def.OrderConsumerService = (*service)(nil)

type service struct {
	orderConsumer             kafka.Consumer
	orderRepository           repository.OrderRepository
	orderShipAssembledDecoder orderDecoderInterface.OrderAssembledDecoder
}

func NewService(orderComsumer kafka.Consumer, orderRepository repository.OrderRepository, orderShipAssembledDecoder orderDecoderInterface.OrderAssembledDecoder) *service {
	return &service{
		orderConsumer:             orderComsumer,
		orderRepository:           orderRepository,
		orderShipAssembledDecoder: orderShipAssembledDecoder,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order consumer service")
	if err := s.orderConsumer.Consume(ctx, s.OrderHandler); err != nil {
		logger.Error(ctx, "Consume from order.assembled topic error", zap.Error(err))
		return err
	}
	return nil
}
