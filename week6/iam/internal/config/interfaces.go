package config

import "time"

type LoggerConfig interface {
	AsJson() bool
	Level() string
}

type IAMGRPCConfig interface {
	Address() string
	Port() string
}

type PostgresConfig interface {
	URL() string
	MigrationDir() string
	DatabaseName() string
}

type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
	CacheTTL() time.Duration
}
