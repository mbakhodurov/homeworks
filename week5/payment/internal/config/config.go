package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/mbakhodurov/homeworks/week5/payment/internal/config/env"
)

var appConfig *config

type config struct {
	Logger      LoggerConfig
	PaymentGRPC PaymentGRPCconfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerConfig, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	paymentGRPCConfig, err := env.NewPaymenyGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:      loggerConfig,
		PaymentGRPC: paymentGRPCConfig,
	}
	return nil
}

func AppConfig() *config {
	return appConfig
}
