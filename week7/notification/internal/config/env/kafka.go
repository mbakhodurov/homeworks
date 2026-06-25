package env

import (
	"github.com/caarlos0/env/v11"
	// def "github.com/mbakhodurov/homeworks/week5/notification/internal/config"
)

type kafkaEnvConfig struct {
	Brokers []string `env:"KAFKA_BROKERS,required"`
}

// var _ def.KafkaConfig = (*kafkaConfig)(nil)

type kafkaConfig struct {
	raw kafkaEnvConfig
}

func NewKafkaConfig() (*kafkaConfig, error) {
	var raw kafkaEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &kafkaConfig{
		raw: raw,
	}, nil
}

func (cfg *kafkaConfig) Brokers() []string {
	return cfg.raw.Brokers
}
