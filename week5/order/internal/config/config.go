package config

import (
	"github.com/joho/godotenv"
	"github.com/mbakhodurov/homeworks/week5/order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger              LoggerConfig
	InventoryGRPCClient InventoryGRPCClientCONFIG
	PaymentGRPCClient   PaymentGRPCClientCONFIG
	OrderHTTP           OrderHTTPConfig
	Postges             PostgesConfig

	Kafka                   KafkaConfig
	OrderPaidProducerConfig OrderPaidProducerConfig

	OrderAssembledConsumerConfig OrderAssembledConsumerConfig
}

func Load(path ...string) error {
	if err := godotenv.Load(path...); err != nil {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	inventoryClientCfg, err := env.NewInventoryClienConfig()
	if err != nil {
		return err
	}

	paymentClientCfg, err := env.NewPaymentClientConfig()
	if err != nil {
		return err
	}

	orderHTTPCfg, err := env.NewOrderHTTPConfig()
	if err != nil {
		return err
	}

	postgesCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderPaidProducerConfig, err := env.NewOrderPaidProducerConfig()
	if err != nil {
		return err
	}

	orderAssembledConsumerConfig, err := env.NewOrderAssembledConsumerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                       loggerCfg,
		InventoryGRPCClient:          inventoryClientCfg,
		PaymentGRPCClient:            paymentClientCfg,
		OrderHTTP:                    orderHTTPCfg,
		Postges:                      postgesCfg,
		Kafka:                        kafkaCfg,
		OrderPaidProducerConfig:      orderPaidProducerConfig,
		OrderAssembledConsumerConfig: orderAssembledConsumerConfig,
	}

	return nil
}

func ApppConfig() *config {
	return appConfig
}
