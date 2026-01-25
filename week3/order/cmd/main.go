package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	v1 "github.com/mbakhodurov/homeworks/week3/order/internal/api/order/v1"
	inventoryGRPCClient "github.com/mbakhodurov/homeworks/week3/order/internal/client/grpc/inventory/v1"
	paymentGRPCClient "github.com/mbakhodurov/homeworks/week3/order/internal/client/grpc/payment/v1"
	"github.com/mbakhodurov/homeworks/week3/order/internal/migrator"
	"github.com/mbakhodurov/homeworks/week3/order/internal/repository/order"
	service "github.com/mbakhodurov/homeworks/week3/order/internal/service/order"
	order_v1 "github.com/mbakhodurov/homeworks/week3/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/mbakhodurov/homeworks/week3/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/mbakhodurov/homeworks/week3/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	inventoryAddress = "localhost:50052"
	paymentAddress   = "localhost:50051"

	httpPort   = "8086"
	configPath = "../../deploy/compose/order/.env"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	ctx := context.Background()

	err := godotenv.Load(configPath)
	if err != nil {
		log.Printf("failed to load .env file: %v\n", err)
		return
	}
	dbURI := os.Getenv("DB_URI")

	// Создаем соединение с базой данных
	con, err := pgx.Connect(ctx, dbURI)
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer func() {
		cerr := con.Close(ctx)
		if cerr != nil {
			log.Printf("failed to close connection: %v\n", cerr)
		}
	}()

	err = con.Ping(ctx)
	if err != nil {
		log.Printf("База данных недоступна: %v\n", err)
		return
	}

	// Инициализируем мигратор
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*con.Config().Copy()), "../"+migrationsDir)

	err = migratorRunner.Up()
	if err != nil {
		log.Printf("Ошибка миграции базы данных: %v\n", err)
		return
	}

	inventoryConn, err := grpc.NewClient(
		inventoryAddress, grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	paymentConn, err := grpc.NewClient(
		paymentAddress, grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	inventoryGRPC := inventory_v1.NewInventoryServiceClient(inventoryConn)
	inventoryClient := inventoryGRPCClient.NewClient(inventoryGRPC)

	paymentGRPC := payment_v1.NewPaymentServiceClient(paymentConn)
	paymentClient := paymentGRPCClient.NewClient(paymentGRPC)

	sqlDB := stdlib.OpenDB(*con.Config().Copy())

	repo := order.NewOrderRepository(sqlDB)
	service := service.NewService(repo, inventoryClient, paymentClient)

	orderHandler := v1.NewApi(service)

	orderServer, err := order_v1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
