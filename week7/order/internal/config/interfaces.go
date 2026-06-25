package config

import (
	"time"

	"github.com/IBM/sarama"
)

type MetricsConfig interface {
	// CollectorEndpoint возвращает адрес OTLP коллектора
	// Обычно: "localhost:4317" (gRPC) или "localhost:4318" (HTTP)
	CollectorEndpoint() string

	// CollectorInterval возвращает интервал отправки метрик
	// Компромисс между актуальностью метрик и нагрузкой на сеть
	// Обычно: 10-60 секунд в зависимости от требований
	CollectorInterval() time.Duration
}

type LoggerConfig interface {
	AsJSON() bool
	Level() string
	EnableOTLP() bool
	OtlpEndpoint() string
	ServiceName() string
	ServiceEnv() string
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
