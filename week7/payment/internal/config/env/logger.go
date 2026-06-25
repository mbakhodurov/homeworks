package env

import "github.com/caarlos0/env/v11"

type loggerEnvConfig struct {
	Level        string `env:"LOGGER_LEVEL,required"`
	AsJson       bool   `env:"LOGGER_AS_JSON,required"`
	EnableOTLP   bool   `env:"LOGGER_ENABLE_OTLP" envDefault:"true"`
	OtlpEndpoint string `env:"OTLP_ENDPOINT" envDefault:"localhost:4317"`
	ServiceName  string `env:"SERVICE_NAME" envDefault:"payment-service"`
	ServiceEnv   string `env:"SERVICE_ENV" envDefault:"dev"`
}

type loggerConfig struct {
	raw loggerEnvConfig
}

func NewLoggerConfig(path ...string) (*loggerConfig, error) {
	var raw loggerEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &loggerConfig{
		raw: raw,
	}, nil
}

func (c *loggerConfig) AsJSON() bool {
	return c.raw.AsJson
}

func (c *loggerConfig) Level() string {
	return c.raw.Level
}

func (c *loggerConfig) EnableOTLP() bool {
	return c.raw.EnableOTLP
}

func (c *loggerConfig) OtlpEndpoint() string {
	return c.raw.OtlpEndpoint
}

func (c *loggerConfig) ServiceName() string {
	return c.raw.ServiceName
}

func (c *loggerConfig) ServiceEnv() string {
	return c.raw.ServiceEnv
}
