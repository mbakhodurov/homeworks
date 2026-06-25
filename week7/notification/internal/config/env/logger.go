package env

import (
	"github.com/caarlos0/env/v11"
	// def "github.com/mbakhodurov/homeworks/week5/notification/internal/config"
)

type loggerEnvConfig struct {
	Level        string `env:"LOGGER_LEVEL,required"`
	AsJson       bool   `env:"LOGGER_AS_JSON,required"`
	EnableOTLP   bool   `env:"LOGGER_ENABLE_OTLP" envDefault:"false"`
	OtlpEndpoint string `env:"OTLP_ENDPOINT" envDefault:"localhost:4317"`
	ServiceName  string `env:"SERVICE_NAME" envDefault:"notification-service"`
	ServiceEnv   string `env:"SERVICE_ENV" envDefault:"dev"`
}

// var _ def.LoggerConfig = (*loggerConfig)(nil

type loggerConfig struct {
	raw loggerEnvConfig
}

func NewLoggerConfig() (*loggerConfig, error) {
	var raw loggerEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &loggerConfig{
		raw: raw,
	}, nil
}

func (cfg *loggerConfig) Level() string {
	return cfg.raw.Level
}

func (cfg *loggerConfig) AsJson() bool {
	return cfg.raw.AsJson
}

func (cfg *loggerConfig) EnableOTLP() bool {
	return cfg.raw.EnableOTLP
}

func (cfg *loggerConfig) OtlpEndpoint() string {
	return cfg.raw.OtlpEndpoint
}

func (cfg *loggerConfig) ServiceName() string {
	return cfg.raw.ServiceName
}

func (cfg *loggerConfig) ServiceEnv() string {
	return cfg.raw.ServiceEnv
}
