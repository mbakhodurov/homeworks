package config

import "time"

type LoggerConfig interface {
	AsJson() bool
	Level() string
}

type InventoryGRPCClientCONFIG interface {
	Address() string
}

type PaymentGRPCClientCONFIG interface {
	Address() string
}

type OrderHTTPConfig interface {
	Address() string
	Readtimeout() time.Duration
	Shutdowntimeout() time.Duration
}

type PostgesConfig interface {
	URL() string
	MigrationDir() string
	DatabaseName() string
}
