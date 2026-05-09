package config

type LoggerConfig interface {
	AsJson() bool
	Level() string
}

type PaymentGRPCconfig interface {
	Address() string
	Port() string
}

type IAMGRPCClientConfig interface {
	Address() string
}
