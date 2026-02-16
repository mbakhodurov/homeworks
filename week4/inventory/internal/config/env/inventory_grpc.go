package env

import "github.com/caarlos0/env"

type inventoryGRPCEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type inventoryGRPCConfig struct {
	raw inventoryGRPCEnvConfig
}

func NewInventoryConfig() (*inventoryGRPCConfig, error) {
	var raw inventoryGRPCEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &inventoryGRPCConfig{raw: raw}, nil
}

func (cfg *inventoryGRPCConfig) Address() string {
	return cfg.raw.Host + ":" + cfg.raw.Port
}
func (cfg *inventoryGRPCConfig) Port() string {
	return cfg.raw.Port
}
