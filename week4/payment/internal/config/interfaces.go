package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type PaymentGRPCconfig interface {
	Address() string
}
