package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/mbakhodurov/homeworks/week7/assembly/internal/config"
	kafkaConverter "github.com/mbakhodurov/homeworks/week7/assembly/internal/converter/kafka"
	"github.com/mbakhodurov/homeworks/week7/assembly/internal/converter/kafka/decoder"
	"github.com/mbakhodurov/homeworks/week7/assembly/internal/service"
	orderproducer "github.com/mbakhodurov/homeworks/week7/assembly/internal/service/producer/order_producer"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/closer"
	wrappedKafka "github.com/mbakhodurov/homeworks/week7/platform/pkg/kafka"
	wrappedKafkaProducer "github.com/mbakhodurov/homeworks/week7/platform/pkg/kafka/producer"
	kafkaMiddleware "github.com/mbakhodurov/homeworks/week7/platform/pkg/middleware/kafka"

	orderConsumer "github.com/mbakhodurov/homeworks/week7/assembly/internal/service/consumer/order_consumer"
	wrappedKafkaConsumer "github.com/mbakhodurov/homeworks/week7/platform/pkg/kafka/consumer"

	"github.com/mbakhodurov/homeworks/week7/platform/pkg/logger"
)

type diContainer struct {
	orderProducerService  service.OrderProducerService
	orderAssemblyProducer wrappedKafka.Producer
	syncProducer          sarama.SyncProducer

	orderConsumerService service.ConsumerService
	orderPaidConsumer    wrappedKafka.Consumer
	consumerGroup        sarama.ConsumerGroup

	orderPaidDecoder kafkaConverter.OrderPaidDecoder
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (di *diContainer) KafkaConverter() kafkaConverter.OrderPaidDecoder {
	if di.orderPaidDecoder == nil {
		di.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}
	return di.orderPaidDecoder
}

func (di *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if di.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return di.consumerGroup.Close()
		})

		di.consumerGroup = consumerGroup
	}
	return di.consumerGroup
}

func (di *diContainer) OrderConsumer() wrappedKafka.Consumer {
	if di.orderPaidConsumer == nil {
		di.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			di.ConsumerGroup(),
			[]string{
				config.AppConfig().OrderPaidConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}
	return di.orderPaidConsumer
}

func (di *diContainer) OrderConsumerService() service.ConsumerService {
	if di.orderConsumerService == nil {
		di.orderConsumerService = orderConsumer.NewService(
			di.OrderConsumer(),
			di.KafkaConverter(),
			di.OrderProducerService(),
		)
	}
	return di.orderConsumerService
}

func (di *diContainer) SyncProducer() sarama.SyncProducer {
	if di.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledProducerConfig.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err))
		}
		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		di.syncProducer = p
	}
	return di.syncProducer
}

func (di *diContainer) OrderAssemblyProducer() wrappedKafka.Producer {
	if di.orderAssemblyProducer == nil {
		di.orderAssemblyProducer = wrappedKafkaProducer.NewProducer(
			di.SyncProducer(),
			config.AppConfig().OrderAssembledProducerConfig.Topic(),
			logger.Logger(),
		)
	}
	return di.orderAssemblyProducer
}

func (di *diContainer) OrderProducerService() service.OrderProducerService {
	if di.orderProducerService == nil {
		di.orderProducerService = orderproducer.NewService(di.OrderAssemblyProducer())
	}
	return di.orderProducerService
}
