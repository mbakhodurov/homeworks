package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/mbakhodurov/homeworks/week4/inventory/internal/app"
	"github.com/mbakhodurov/homeworks/week4/inventory/internal/config"
	"github.com/mbakhodurov/homeworks/week4/platform/pkg/closer"
	"github.com/mbakhodurov/homeworks/week4/platform/pkg/logger"
	"go.uber.org/zap"
)

const configPath = "../../deploy/compose/inventory/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.NewApp(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Не удалось создать приложение", zap.Error(err))
		return
	}

	err = a.Run(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Ошибка при работе приложения", zap.Error(err))
		return
	}

	// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// 	defer cancel()

	// 	// Создаем MongoDB клиент
	// 	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
	// 	if err != nil {
	// 		log.Printf("failed to connect to MongoDB: %v\n", err)
	// 		return
	// 	}
	// 	defer func() {
	// 		if cerr := mongoClient.Disconnect(context.Background()); cerr != nil {
	// 			log.Printf("failed to disconnect from MongoDB: %v\n", cerr)
	// 		}
	// 	}()

	// 	// Проверяем подключение к MongoDB
	// 	err = mongoClient.Ping(ctx, nil)
	// 	if err != nil {
	// 		log.Printf("failed to ping MongoDB: %v\n", err)
	// 		return
	// 	}
	// 	log.Println("✅ Connected to MongoDB")

	// 	lis, err := net.Listen("tcp", config.AppConfig().InventoryGRPC.Address())
	// 	if err != nil {
	// 		log.Printf("failed to listen: %v\n", err)
	// 		return
	// 	}
	// 	defer func() {
	// 		if cerr := lis.Close(); cerr != nil {
	// 			log.Printf("failed to close listener: %v\n", cerr)
	// 		}
	// 	}()

	// 	// Создаем gRPC сервер
	// 	s := grpc.NewServer(
	// 		grpc.ChainUnaryInterceptor(
	// 			grpc.UnaryServerInterceptor(interceptor.LoggerInterceptor()),
	// 			// grpc.UnaryServerInterceptor(interceptor.LoggerInterceptor2()),
	// 		),
	// 	)

	// 	logger.Init(config.AppConfig().Logger.Level(),
	// 		config.AppConfig().Logger.AsJson(),
	// 	)

	// 	db := mongoClient.Database("inventory")
	// 	repo := repo.NewRepository(ctx, db)
	// 	service := service.NewInventoryService(repo)
	// 	api := inventoryV1.NewInventoryApi(service)

	// 	inventory_v1.RegisterInventoryServiceServer(s, api)

	// 	// Включаем рефлексию для отладки
	// 	reflection.Register(s)

	// 	go func() {
	// 		log.Printf("🚀 gRPC server listening on %s\n", config.AppConfig().InventoryGRPC.Address())
	// 		err = s.Serve(lis)
	// 		if err != nil {
	// 			log.Printf("failed to serve: %v\n", err)
	// 			return
	// 		}
	// 	}()

	// 	// Graceful shutdown
	// 	quit := make(chan os.Signal, 1)
	// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 	<-quit
	// 	log.Println("🛑 Shutting down gRPC server...")
	// 	s.GracefulStop()
	// 	log.Println("✅ Server stopped")

}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы", zap.Error(err))
	}
}
