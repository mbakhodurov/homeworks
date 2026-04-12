package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
	// def "github.com/mbakhodurov/homeworks/week5/notification/internal/config"
)

// var _ def.OrderPaidConsumerConfig = (*orderPaidConsumerConfig)(nil)

type orderPaidConsumerEnvConfig struct {
	TopicName string `env:"ORDER_PAID_TOPIC_NAME,required"`
	GroupID   string `env:"ORDER_PAID_CONSUMER_GROUP_ID,required"`
}

type orderPaidConsumerConfig struct {
	raw orderPaidConsumerEnvConfig
}

func NewOrderPaidConsumerConfig() (*orderPaidConsumerConfig, error) {
	var raw orderPaidConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderPaidConsumerConfig{raw: raw}, nil
}

func (cfg *orderPaidConsumerConfig) GroupID() string {
	return cfg.raw.GroupID
}

func (cfg *orderPaidConsumerConfig) Topic() string {
	return cfg.raw.TopicName
}

func (cfg *orderPaidConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true

	return config
}
