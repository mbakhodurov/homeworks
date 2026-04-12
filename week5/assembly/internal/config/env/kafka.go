package env

import "github.com/caarlos0/env/v11"

type KafkaEnvConfig struct {
	Brokers []string `env:"KAFKA_BROKERS,required"`
}

type kafkaConfig struct {
	raw KafkaEnvConfig
}

func NewKafkaConfig() (*kafkaConfig, error) {
	var raw KafkaEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &kafkaConfig{raw: raw}, nil
}

func (cfg *kafkaConfig) Brokers() []string {
	return cfg.raw.Brokers
}
