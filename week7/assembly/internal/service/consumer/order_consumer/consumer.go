package order_consumer

import (
	"context"

	kafkaConverter "github.com/mbakhodurov/homeworks/week7/assembly/internal/converter/kafka"
	def "github.com/mbakhodurov/homeworks/week7/assembly/internal/service"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"

	"github.com/mbakhodurov/homeworks/week7/platform/pkg/kafka"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/logger"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	orderConsumer    kafka.Consumer
	orderPaidDecoder kafkaConverter.OrderPaidDecoder
	orderProducer    def.OrderProducerService

	assemblyHistogram metric.Float64Histogram
}

func NewService(orderConsumer kafka.Consumer, orderPaidDecoder kafkaConverter.OrderPaidDecoder, orderProducer def.OrderProducerService) *service {
	meter := otel.Meter("assembly-service")

	assemblyHistogram, _ := meter.Float64Histogram(
		// "assembly_duration_seconds", // OTel добавляет _seconds (unit=s) автоматически → получалось assembly_duration_seconds_seconds
		"assembly_duration",
		metric.WithDescription("Duration of order assembly in seconds"),
		metric.WithUnit("s"),
	)

	return &service{
		orderConsumer:     orderConsumer,
		orderPaidDecoder:  orderPaidDecoder,
		orderProducer:     orderProducer,
		assemblyHistogram: assemblyHistogram,
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
