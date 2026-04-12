package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/go-telegram/bot"
	"github.com/mbakhodurov/homeworks/week5/notification/internal/client/http"
	httpClient "github.com/mbakhodurov/homeworks/week5/notification/internal/client/http"
	telegramClient "github.com/mbakhodurov/homeworks/week5/notification/internal/client/http/telegram"
	"github.com/mbakhodurov/homeworks/week5/notification/internal/config"
	"github.com/mbakhodurov/homeworks/week5/notification/internal/converter/kafka"
	kafkaConverter "github.com/mbakhodurov/homeworks/week5/notification/internal/converter/kafka"
	"github.com/mbakhodurov/homeworks/week5/notification/internal/converter/kafka/decoder"
	"github.com/mbakhodurov/homeworks/week5/notification/internal/service"
	"github.com/mbakhodurov/homeworks/week5/notification/internal/service/consumer/order_assembled_consumer"
	"github.com/mbakhodurov/homeworks/week5/notification/internal/service/consumer/order_paid_consumer"
	"github.com/mbakhodurov/homeworks/week5/notification/internal/service/telegram"

	wrappedKafka "github.com/mbakhodurov/homeworks/week5/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/mbakhodurov/homeworks/week5/platform/pkg/kafka/consumer"
	"github.com/mbakhodurov/homeworks/week5/platform/pkg/logger"
	kafkaMiddleware "github.com/mbakhodurov/homeworks/week5/platform/pkg/middleware/kafka"
)

type diContainer struct {
	telegramClient           http.TelegramClient
	orderPaidDecoder         kafka.OrderPaidDecoder
	orderPaidConsumerGroup   sarama.ConsumerGroup
	orderPaidConsumerService service.ConsumerService
	orderPaidConsumer        wrappedKafka.Consumer
	telegramService          service.TelegramService
	telegramBot              *bot.Bot

	orderAssembledDecoder         kafka.OrderAssembledDecoder
	orderAssembledConsumerService service.ConsumerService
	orderAssembledConsumer        wrappedKafka.Consumer
	orderAssembledConsumerGroup   sarama.ConsumerGroup
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderAssembledConsumerGroup() sarama.ConsumerGroup {
	if d.orderAssembledConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.GroupID(),
			config.AppConfig().OrderAssembledConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create order assembled consumer group: %s\n", err))
		}

		d.orderAssembledConsumerGroup = consumerGroup
	}
	return d.orderAssembledConsumerGroup
}

func (di *diContainer) OrderAssembledConsumer() wrappedKafka.Consumer {
	if di.orderAssembledConsumer == nil {
		di.orderAssembledConsumer = wrappedKafkaConsumer.NewConsumer(
			di.OrderAssembledConsumerGroup(),
			[]string{
				config.AppConfig().OrderAssembledConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}
	return di.orderAssembledConsumer

}

func (di *diContainer) OrderAssembledConsumerService(ctx context.Context) service.ConsumerService {
	if di.orderAssembledConsumerService == nil {
		di.orderAssembledConsumerService = order_assembled_consumer.NewService(
			di.OrderAssembledConsumer(),
			di.OrderAssembledDecoder(),
			di.TelegramService(ctx),
		)
	}
	return di.orderAssembledConsumerService
}

func (di *diContainer) OrderAssembledDecoder() kafka.OrderAssembledDecoder {
	if di.orderAssembledDecoder == nil {
		di.orderAssembledDecoder = decoder.NewOrderAssembledDecoder()
	}
	return di.orderAssembledDecoder
}

func (di *diContainer) TelegramBot(ctx context.Context) *bot.Bot {
	if di.telegramBot == nil {
		b, err := bot.New(
			config.AppConfig().TelegramBot.Token(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create telegram bot: %s\n", err))
		}
		di.telegramBot = b
	}
	return di.telegramBot
}

func (di *diContainer) TelegramClient(ctx context.Context) httpClient.TelegramClient {
	if di.telegramClient == nil {
		di.telegramClient = telegramClient.NewClient(di.TelegramBot(ctx))
	}
	return di.telegramClient
}

func (di *diContainer) TelegramService(ctx context.Context) service.TelegramService {
	if di.telegramService == nil {
		di.telegramService = telegram.NewService(di.TelegramClient(ctx))
	}
	return di.telegramService
}

func (di *diContainer) OrderPaidDecoder() kafkaConverter.OrderPaidDecoder {
	if di.orderPaidDecoder == nil {
		di.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}
	return di.orderPaidDecoder
}

func (di *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if di.orderPaidConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create order paid consumer group: %s\n", err))
		}

		di.orderPaidConsumerGroup = consumerGroup
	}
	return di.orderPaidConsumerGroup
}

func (di *diContainer) OrderPaidConsumer() wrappedKafka.Consumer {
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

func (di *diContainer) OrderPaidConsumerService(ctx context.Context) service.ConsumerService {
	if di.orderPaidConsumerService == nil {
		di.orderPaidConsumerService = order_paid_consumer.NewService(
			di.OrderPaidConsumer(),
			di.OrderPaidDecoder(),
			di.TelegramService(ctx),
		)
	}
	return di.orderPaidConsumerService
}
