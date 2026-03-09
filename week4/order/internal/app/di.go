package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"

	orderV1API "github.com/mbakhodurov/homeworks/week4/order/internal/api/order/v1"
	client "github.com/mbakhodurov/homeworks/week4/order/internal/client/grpc"
	inventortyClientV1 "github.com/mbakhodurov/homeworks/week4/order/internal/client/grpc/inventory/v1"
	"github.com/mbakhodurov/homeworks/week4/order/internal/config"
	"github.com/mbakhodurov/homeworks/week4/order/internal/migrator"
	repoInterface "github.com/mbakhodurov/homeworks/week4/order/internal/repository"
	repo "github.com/mbakhodurov/homeworks/week4/order/internal/repository/order"
	serviceInterface "github.com/mbakhodurov/homeworks/week4/order/internal/service"
	service "github.com/mbakhodurov/homeworks/week4/order/internal/service/order"
	"github.com/mbakhodurov/homeworks/week4/platform/pkg/closer"
	order_v1 "github.com/mbakhodurov/homeworks/week4/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/mbakhodurov/homeworks/week4/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/mbakhodurov/homeworks/week4/shared/pkg/proto/payment/v1"

	// payment_v1 "github.com/mbakhodurov/homeworks/week4/shared/pkg/proto/payment/v1"
	paymentClientV1 "github.com/mbakhodurov/homeworks/week4/order/internal/client/grpc/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type diContainer struct {
	pgxConn *pgx.Conn
	sqlDb   *sql.DB

	migrationRunner *migrator.Migrator

	orderRepository repoInterface.OrderRepository

	inventoryGRPCConn *grpc.ClientConn
	paymentGRPCConn   *grpc.ClientConn

	inventoryClient client.InventoryClient
	paymenyClient   client.PaymentClient

	orderService  serviceInterface.OrderService
	orderV1API    order_v1.Handler
	orderV1Server *order_v1.Server
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (di *diContainer) OrderV1Server(ctx context.Context) *order_v1.Server {
	if di.orderV1Server == nil {
		server, err := order_v1.NewServer(di.OrderV1API(ctx))
		if err != nil {
			panic(fmt.Sprintf("💥 failed to create OrderV1 server: %v", err))
		}
		di.orderV1Server = server
	}
	return di.orderV1Server
}

func (di *diContainer) OrderV1API(ctx context.Context) order_v1.Handler {
	if di.orderV1API == nil {
		di.orderV1API = orderV1API.NewApi(di.OrderService(ctx))
	}

	return di.orderV1API
}

func (di *diContainer) OrderService(ctx context.Context) serviceInterface.OrderService {
	if di.orderService == nil {
		di.orderService = service.NewService(di.OrderRepository(ctx), di.InventoryClient(ctx), di.PaymentClient(ctx))
	}

	return di.orderService
}

func (di *diContainer) PaymentClient(ctx context.Context) client.PaymentClient {
	if di.paymenyClient == nil {
		di.paymenyClient = paymentClientV1.NewClient(payment_v1.NewPaymentServiceClient(di.PaymentGRPCConn(ctx)))
	}

	return di.paymenyClient
}

func (di *diContainer) InventoryClient(ctx context.Context) client.InventoryClient {
	if di.inventoryClient == nil {
		di.inventoryClient = inventortyClientV1.NewClient(inventory_v1.NewInventoryServiceClient(di.InventoryGRPCConn(ctx)))
	}

	return di.inventoryClient
}

func (di *diContainer) PaymentGRPCConn(ctx context.Context) *grpc.ClientConn {
	if di.paymentGRPCConn == nil {
		conn, err := grpc.NewClient(
			config.ApppConfig().PaymentGRPCClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("💥 failed to connect to payment service: %v", err))
		}

		closer.AddNamed("Payment gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		di.paymentGRPCConn = conn
	}

	return di.paymentGRPCConn
}

func (di *diContainer) InventoryGRPCConn(ctx context.Context) *grpc.ClientConn {
	if di.inventoryGRPCConn == nil {
		conn, err := grpc.NewClient(
			config.ApppConfig().InventoryGRPCClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("💥 failed to connect to inventory service: %v", err))
		}

		closer.AddNamed("Inventory gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		di.inventoryGRPCConn = conn
	}

	return di.inventoryGRPCConn
}

func (di *diContainer) OrderRepository(ctx context.Context) repoInterface.OrderRepository {
	if di.orderRepository == nil {
		di.orderRepository = repo.NewRepository(di.SqlDB(ctx))
	}

	return di.orderRepository
}

func (d *diContainer) MigratorRunner(ctx context.Context) *migrator.Migrator {
	if d.migrationRunner == nil {
		d.migrationRunner = migrator.NewMigrator(d.SqlDB(ctx), config.ApppConfig().Postges.MigrationDir())
	}
	return d.migrationRunner
}

func (di *diContainer) SqlDB(ctx context.Context) *sql.DB {
	if di.sqlDb == nil {
		di.sqlDb = stdlib.OpenDB(*di.PgxConn(ctx).Config().Copy())
	}
	return di.sqlDb
}

func (di *diContainer) PgxConn(ctx context.Context) *pgx.Conn {
	if di.pgxConn == nil {
		conn, err := pgx.Connect(ctx, config.ApppConfig().Postges.URL())
		if err != nil {
			panic(fmt.Sprintf("💥 failed to connect to database: %v", err))
		}

		err = conn.Ping(ctx)
		if err != nil {
			panic(fmt.Sprintf("💥 failed to ping database: %v", err))
		}

		closer.AddNamed("PostgresSQL connection", func(ctx context.Context) error {
			return conn.Close(ctx)
		})

		di.pgxConn = conn
	}
	return di.pgxConn
}
