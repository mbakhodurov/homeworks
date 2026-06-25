package app

import (
	"context"
	"fmt"

	v1 "github.com/mbakhodurov/homeworks/week7/inventory/internal/api/inventory/v1"
	"github.com/mbakhodurov/homeworks/week7/inventory/internal/config"
	"github.com/mbakhodurov/homeworks/week7/inventory/internal/repository"
	inventoryRepo "github.com/mbakhodurov/homeworks/week7/inventory/internal/repository/inventory"
	"github.com/mbakhodurov/homeworks/week7/inventory/internal/service"
	inventoryService "github.com/mbakhodurov/homeworks/week7/inventory/internal/service/inventory"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/closer"
	auth_v1 "github.com/mbakhodurov/homeworks/week7/shared/pkg/proto/auth/v1"
	inventory_v1 "github.com/mbakhodurov/homeworks/week7/shared/pkg/proto/inventory/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type diContainer struct {
	inventoryRepo    repository.InventoryRepository
	inventoryService service.InventoryService
	inventoryV1Api   inventory_v1.InventoryServiceServer

	mongoDBHandle *mongo.Database
	mongoDBClient *mongo.Client

	iamGRPCConn *grpc.ClientConn
	iamClient   auth_v1.AuthServiceClient
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) IAMClient(ctx context.Context) auth_v1.AuthServiceClient {
	if d.iamClient == nil {
		d.iamClient = auth_v1.NewAuthServiceClient(d.IAMGRPCConn(ctx))
	}
	return d.iamClient
}

func (d *diContainer) IAMGRPCConn(_ context.Context) *grpc.ClientConn {
	if d.iamGRPCConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().IAMClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("💥 failed to connect to IAM service: %v", err))
		}

		closer.AddNamed("IAM gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.iamGRPCConn = conn
	}
	return d.iamGRPCConn
}

func (di *diContainer) InventoryV1Api(ctx context.Context) inventory_v1.InventoryServiceServer {
	if di.inventoryV1Api == nil {
		di.inventoryV1Api = v1.NewInventoryV1Api(di.InventoryService(ctx))
	}
	return di.inventoryV1Api
}

func (di *diContainer) InventoryService(ctx context.Context) service.InventoryService {
	if di.inventoryService == nil {
		di.inventoryService = inventoryService.NewInventoryService(di.InventoryRepo(ctx))
	}
	return di.inventoryService
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

func (di *diContainer) InventoryRepo(ctx context.Context) repository.InventoryRepository {
	if di.inventoryRepo == nil {
		di.inventoryRepo = inventoryRepo.NewRepository(ctx, di.MongoDBHandle(ctx))
	}
	return di.inventoryRepo
}
