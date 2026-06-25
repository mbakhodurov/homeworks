package config

type LoggerConfig interface {
	AsJSON() bool
	Level() string
	EnableOTLP() bool
	OtlpEndpoint() string
	ServiceName() string
	ServiceEnv() string
}

type PaymentGRPCconfig interface {
	Address() string
	Port() string
}

type IAMGRPCClientConfig interface {
	Address() string
}
