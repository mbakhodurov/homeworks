package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type paymentGRPEnvCconfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type paymentGRPCconfig struct {
	raw paymentGRPEnvCconfig
}

func NewPaymenyGRPCConfig() (*paymentGRPCconfig, error) {
	var raw paymentGRPEnvCconfig

	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &paymentGRPCconfig{
		raw: raw,
	}, nil
}

func (cfg *paymentGRPCconfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
