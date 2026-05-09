package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/mbakhodurov/homeworks/week6/inventory/internal/config/env"
)

var appConfig *config

type config struct {
	Logger        LoggerConfig
	Mongo         MongoConfig
	InventoryGRPC InventoryGRPCConfig
	IAMClient     IAMGRPCClientConfig
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

	mongoConfig, err := env.NewMongoConfig()
	if err != nil {
		return err
	}

	inventoryGRPCConfig, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	iamGRPCConfig, err := env.NewIAMGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:        loggerConfig,
		Mongo:         mongoConfig,
		InventoryGRPC: inventoryGRPCConfig,
		IAMClient:     iamGRPCConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
