package env

import (
	"github.com/caarlos0/env/v11"
)

type paymentClientEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type paymentClientConfig struct {
	raw paymentClientEnvConfig
}

func NewPaymentClientConfig() (*paymentClientConfig, error) {
	var raw paymentClientEnvConfig
	if err := env.Parse(raw); err != nil {
		return nil, err
	}

	return &paymentClientConfig{
		raw: raw,
	}, nil
}

func (cfg *paymentClientConfig) Address() string {
	return cfg.raw.Host + ":" + cfg.raw.Port
}
