package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/mbakhodurov/homeworks/week7/order/internal/app"
	"github.com/mbakhodurov/homeworks/week7/order/internal/config"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/closer"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/logger"
	"go.uber.org/zap"
)

const (
	configPath = "../../deploy/compose/order/.env"

// 	inventoryAddress = "localhost:50052"
// 	paymentAddress   = "localhost:50051"

// httpPort = "8086"
// // Таймауты для HTTP-сервера
// readHeaderTimeout = 5 * time.Second
// shutdownTimeout   = 10 * time.Second
)

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
	// ctx := context.Background()

	// err := config.Load(configPath)
	// if err != nil {
	// 	panic(fmt.Errorf("failed to load config: %w", err))
	// }

	// fmt.Println("config", config.ApppConfig().Logger.Level())

	// logger.Init(config.ApppConfig().Logger.Level(), config.ApppConfig().Logger.AsJson())

	// // fmt.Println("Config:", config.ApppConfig().InventoryGRPCClient.Address())

	// con, err := pgx.Connect(ctx, config.ApppConfig().Postges.URL())
	// if err != nil {
	// 	log.Printf("failed to connect to database: %v\n", err)
	// 	return
	// }

	// defer func() {
	// 	cerr := con.Close(ctx)
	// 	if cerr != nil {
	// 		log.Printf("failed to close connection: %v\n", cerr)
	// 	}
	// }()

	// err = con.Ping(ctx)
	// if err != nil {
	// 	log.Printf("База данных недоступна: %v\n", err)
	// 	return
	// }

	// sqlDB := stdlib.OpenDB(*con.Config().Copy())

	// // Инициализируем мигратор
	// // migrationsDir := config.ApppConfig().Postges.MigrationDir()
	// // migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*con.Config().Copy()), "../."+migrationsDir)
	// // err = migratorRunner.Up(ctx)
	// // if err != nil {
	// // 	log.Printf("Ошибка миграции базы данных: %v\n", err)
	// // 	return
	// // }

	// inventoryConn, err := grpc.NewClient(
	// 	config.ApppConfig().InventoryGRPCClient.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()),
	// )
	// if err != nil {
	// 	log.Printf("failed to connect: %v\n", err)
	// 	return
	// }
	// defer func() {
	// 	if cerr := inventoryConn.Close(); cerr != nil {
	// 		log.Printf("failed to close connect: %v", cerr)
	// 	}
	// }()

	// paymentConn, err := grpc.NewClient(
	// 	config.ApppConfig().InventoryGRPCClient.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()),
	// )
	// if err != nil {
	// 	log.Printf("failed to connect: %v\n", err)
	// 	return
	// }
	// defer func() {
	// 	if cerr := paymentConn.Close(); cerr != nil {
	// 		log.Printf("failed to close connect: %v", cerr)
	// 	}
	// }()

	// inventoryGRPC := inventory_v1.NewInventoryServiceClient(inventoryConn)
	// inventoryClient := inventoryGRPCClient.NewClient(inventoryGRPC)

	// paymentGRPC := payment_v1.NewPaymentServiceClient(paymentConn)
	// paymentClient := paymentGRPCClient.NewClient(paymentGRPC)

	// repo := order.NewRepository(sqlDB)
	// service := service.NewService(repo, inventoryClient, paymentClient)

	// orderHandler := v1.NewApi(service)

	// orderServer, err := order_v1.NewServer(orderHandler)
	// if err != nil {
	// 	log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	// }

	// r := chi.NewRouter()
	// r.Use(middleware.Logger)
	// r.Use(middleware.Recoverer)
	// r.Use(middleware.Timeout(10 * time.Second))

	// r.Mount("/", orderServer)

	// server := &http.Server{
	// 	Addr:              net.JoinHostPort("localhost", httpPort),
	// 	Handler:           r,
	// 	ReadHeaderTimeout: config.ApppConfig().OrderHTTP.Readtimeout(),
	// }

	// go func() {
	// 	log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
	// 	err = server.ListenAndServe()
	// 	if err != nil && !errors.Is(err, http.ErrServerClosed) {
	// 		log.Printf("❌ Ошибка запуска сервера: %v\n", err)
	// 	}
	// }()

	// // Graceful shutdown
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// <-quit

	// log.Println("🛑 Завершение работы сервера...")

	// // Создаем контекст с таймаутом для остановки сервера
	// ctx, cancel := context.WithTimeout(context.Background(), config.ApppConfig().OrderHTTP.Shutdowntimeout())
	// defer cancel()

	// err = server.Shutdown(ctx)
	// if err != nil {
	// 	log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	// }

	// log.Println("✅ Сервер остановлен")
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы", zap.Error(err))
	}
}
