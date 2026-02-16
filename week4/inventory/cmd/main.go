package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	inventoryV1 "github.com/mbakhodurov/homeworks/week4/inventory/internal/api/inventory/v1"
	"github.com/mbakhodurov/homeworks/week4/inventory/internal/config"
	"github.com/mbakhodurov/homeworks/week4/inventory/internal/interceptor"
	repo "github.com/mbakhodurov/homeworks/week4/inventory/internal/repository/inventory"
	service "github.com/mbakhodurov/homeworks/week4/inventory/internal/service/inventory"
	"github.com/mbakhodurov/homeworks/week4/platform/pkg/logger"
	inventory_v1 "github.com/mbakhodurov/homeworks/week4/shared/pkg/proto/inventory/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const configPath = "../../deploy/compose/inventory/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Создаем MongoDB клиент
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
	if err != nil {
		log.Printf("failed to connect to MongoDB: %v\n", err)
		return
	}
	defer func() {
		if cerr := mongoClient.Disconnect(context.Background()); cerr != nil {
			log.Printf("failed to disconnect from MongoDB: %v\n", cerr)
		}
	}()

	// Проверяем подключение к MongoDB
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Printf("failed to ping MongoDB: %v\n", err)
		return
	}
	log.Println("✅ Connected to MongoDB")

	lis, err := net.Listen("tcp", config.AppConfig().InventoryGRPC.Address())
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	// Создаем gRPC сервер
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc.UnaryServerInterceptor(interceptor.LoggerInterceptor()),
			// grpc.UnaryServerInterceptor(interceptor.LoggerInterceptor2()),
		),
	)

	// Регистрируем наш сервис

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
	if err != nil {
		panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
	}
	defer func() {
		cerr := client.Disconnect(ctx)
		if cerr != nil {
			log.Printf("failed to disconnect: %v\n", cerr)
		}
	}()

	logger.Init(config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)

	db := mongoClient.Database("inventory")
	repo := repo.NewRepository(ctx, db)
	service := service.NewInventoryService(repo)
	api := inventoryV1.NewInventoryApi(service)

	inventory_v1.RegisterInventoryServiceServer(s, api)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %s\n", config.AppConfig().InventoryGRPC.Address())
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")

}
