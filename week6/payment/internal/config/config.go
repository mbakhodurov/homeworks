package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/mbakhodurov/homeworks/week6/payment/internal/config/env"
)

var appConfig *config

type config struct {
	Logger      LoggerConfig
	PaymentGRPC PaymentGRPCconfig
	IAMClient   IAMGRPCClientConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	logger, err := env.NewLoggerConfig(path...)
	if err != nil {
		return err
	}

	paymentGRPCConfig, err := env.NewPaymenyGRPCConfig()
	if err != nil {
		return err
	}

	iamGRPCConfig, err := env.NewIAMGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:      logger,
		PaymentGRPC: paymentGRPCConfig,
		IAMClient:   iamGRPCConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
