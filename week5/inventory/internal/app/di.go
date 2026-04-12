package app

import (
	"context"
	"fmt"

	v1 "github.com/mbakhodurov/homeworks/week5/inventory/internal/api/inventory/v1"
	"github.com/mbakhodurov/homeworks/week5/inventory/internal/config"
	"github.com/mbakhodurov/homeworks/week5/inventory/internal/repository"
	"github.com/mbakhodurov/homeworks/week5/inventory/internal/service"
	"github.com/mbakhodurov/homeworks/week5/platform/pkg/closer"
	inventory_v1 "github.com/mbakhodurov/homeworks/week5/shared/pkg/proto/inventory/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	inventoryRepo "github.com/mbakhodurov/homeworks/week5/inventory/internal/repository/inventory"
	inventoryService "github.com/mbakhodurov/homeworks/week5/inventory/internal/service/inventory"
)

type diContainer struct {
	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database

	inventoryRepo    repository.InventoryRepository
	inventoryService service.InventoryService

	inventoryApi inventory_v1.InventoryServiceServer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (di *diContainer) InventoryV1Api(ctx context.Context) inventory_v1.InventoryServiceServer {
	if di.inventoryApi == nil {
		di.inventoryApi = v1.NewInventoryV1Api(di.InventoryService(ctx))
	}
	return di.inventoryApi
}

func (di *diContainer) InventoryService(ctx context.Context) service.InventoryService {
	if di.inventoryService == nil {
		di.inventoryService = inventoryService.NewInventoryService(di.InventoryRepo(ctx))
	}
	return di.inventoryService
}

func (di *diContainer) InventoryRepo(ctx context.Context) repository.InventoryRepository {
	if di.inventoryRepo == nil {
		di.inventoryRepo = inventoryRepo.NewRepository(ctx, di.MongoDBHandle(ctx))
	}
	return di.inventoryRepo
}

func (di *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if di.mongoDBHandle == nil {
		di.mongoDBHandle = di.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}
	return di.mongoDBHandle
}

func (di *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if di.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v\n", err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		di.mongoDBClient = client
	}

	return di.mongoDBClient
}
