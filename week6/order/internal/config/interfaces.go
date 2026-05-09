package config

import (
	"time"

	"github.com/IBM/sarama"
)

type LoggerConfig interface {
	AsJson() bool
	Level() string
}

type InventoryGRPCClientCONFIG interface {
	Address() string
}

type PaymentGRPCClientCONFIG interface {
	Address() string
}

type IAMGRPCClientConfig interface {
	Address() string
}

type OrderHTTPConfig interface {
	Address() string
	Readtimeout() time.Duration
	Shutdowntimeout() time.Duration
}

type PostgesConfig interface {
	URL() string
	MigrationDir() string
	DatabaseName() string
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type OrderAssembledConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}
