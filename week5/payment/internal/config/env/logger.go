package env

import "github.com/caarlos0/env/v11"

type loggerEnvConfig struct {
	Level  string `env:"LOGGER_LEVEL,required"`
	AsJson bool   `env:"LOGGER_AS_JSON,required"`
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

func (l *loggerConfig) Level() string {
	return l.raw.Level
}

func (logger *loggerConfig) AsJson() bool {
	return logger.raw.AsJson
}
