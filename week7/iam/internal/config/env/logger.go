package env

import (
	"github.com/caarlos0/env/v11"
)

type loggerEnvConfig struct {
	Level        string `env:"LOGGER_LEVEL,required"`
	AsJson       bool   `env:"LOGGER_AS_JSON,required"`
	EnableOTLP   bool   `env:"LOGGER_ENABLE_OTLP" envDefault:"false"`
	OtlpEndpoint string `env:"OTLP_ENDPOINT" envDefault:"localhost:4317"`
	ServiceName  string `env:"SERVICE_NAME" envDefault:"iam-service"`
	ServiceEnv   string `env:"SERVICE_ENV" envDefault:"dev"`
}

type loggerConfig struct {
	raw loggerEnvConfig
}

func NewLoggerConfig() (*loggerConfig, error) {
	var raw loggerEnvConfig

	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &loggerConfig{
		raw: raw,
	}, nil
}

func (l *loggerConfig) AsJson() bool {
	return l.raw.AsJson
}

func (cfg *loggerConfig) Level() string {
	return cfg.raw.Level
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
