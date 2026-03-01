package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type orderHTTPEnvConfig struct {
	Host            string        `env:"HTTP_HOST,required"`
	Port            string        `env:"HTTP_PORT,required"`
	ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT,required"`
	ShutdownTimeout time.Duration `env:"HTTP_SHUT_DOWN_TIMEOUT,required"`
}

type orderHTTPConfig struct {
	raw orderHTTPEnvConfig
}

func NewOrderHTTPConfig() (*orderHTTPConfig, error) {
	var raw orderHTTPEnvConfig
	if err := env.Parse(raw); err != nil {
		return nil, err
	}

	return &orderHTTPConfig{
		raw: raw,
	}, nil
}

func (cfg *orderHTTPConfig) Address() string {
	return cfg.raw.Host + ":" + cfg.raw.Port
}

func (cfg *orderHTTPConfig) Readtimeout() time.Duration {
	return cfg.raw.ReadTimeout
}

func (cfg *orderHTTPConfig) Shutdowntimeout() time.Duration {
	return cfg.raw.ShutdownTimeout
}

// Address() string
// Readtimeout() time.Duration
// Shutdowntimeout() time.Duration
