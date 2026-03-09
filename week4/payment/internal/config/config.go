package config

import (
	"github.com/joho/godotenv"
	"github.com/mbakhodurov/homeworks/week4/payment/internal/config/env"
)

var appConfig *config

type config struct {
	Logger      LoggerConfig
	PaymentGRPC PaymentGRPCconfig
}

func Load(path ...string) error {
	if err := godotenv.Load(path...); err != nil {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig(path...)
	if err != nil {
		return err
	}

	paymentGRPCCfg, err := env.NewPaymenyGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		PaymentGRPC: paymentGRPCCfg,
		Logger:      loggerCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
