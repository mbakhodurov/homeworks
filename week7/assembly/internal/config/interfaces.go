package config

import (
	"time"

	"github.com/IBM/sarama"
)

type LoggerConfig interface {
	Level() string
	AsJson() bool
	EnableOTLP() bool
	OtlpEndpoint() string
	ServiceName() string
	ServiceEnv() string
}

type MetricsConfig interface {
	CollectorEndpoint() string
	CollectorInterval() time.Duration
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}

type OrderAssembledProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}
