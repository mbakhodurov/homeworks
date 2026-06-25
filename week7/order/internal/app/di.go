package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	client "github.com/mbakhodurov/homeworks/week7/order/internal/client/grpc"
	inventoryClient "github.com/mbakhodurov/homeworks/week7/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/mbakhodurov/homeworks/week7/order/internal/client/grpc/payment/v1"
	"github.com/mbakhodurov/homeworks/week7/order/internal/converter/kafka"
	"github.com/mbakhodurov/homeworks/week7/order/internal/converter/kafka/decoder"

	serviceHandler "github.com/mbakhodurov/homeworks/week7/order/internal/api/order/v1"
	"github.com/mbakhodurov/homeworks/week7/order/internal/config"
	repoOrderInterface "github.com/mbakhodurov/homeworks/week7/order/internal/repository"
	repoOrder "github.com/mbakhodurov/homeworks/week7/order/internal/repository/order"
	serviceOrderInterface "github.com/mbakhodurov/homeworks/week7/order/internal/service"
	serviceOrder "github.com/mbakhodurov/homeworks/week7/order/internal/service/order"
	orderProducer "github.com/mbakhodurov/homeworks/week7/order/internal/service/producer/order_producer"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/closer"
	wrappedKafka "github.com/mbakhodurov/homeworks/week7/platform/pkg/kafka"
	wrappedKafkaProducer "github.com/mbakhodurov/homeworks/week7/platform/pkg/kafka/producer"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/logger"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/migrator"
	order_v1 "github.com/mbakhodurov/homeworks/week7/shared/pkg/openapi/order/v1"
	auth_v1 "github.com/mbakhodurov/homeworks/week7/shared/pkg/proto/auth/v1"
	inventory_v1 "github.com/mbakhodurov/homeworks/week7/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/mbakhodurov/homeworks/week7/shared/pkg/proto/payment/v1"

	pgMigrator "github.com/mbakhodurov/homeworks/week7/platform/pkg/migrator/pg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderconsumer "github.com/mbakhodurov/homeworks/week7/order/internal/service/consumer/order_consumer"
	wrappedKafkaConsumer "github.com/mbakhodurov/homeworks/week7/platform/pkg/kafka/consumer"
	kafkaMiddleware "github.com/mbakhodurov/homeworks/week7/platform/pkg/middleware/kafka"
)

type diContainer struct {
	syncProducer         sarama.SyncProducer
	orderProducerService serviceOrderInterface.OrderProducerService
	orderProducer        wrappedKafka.Producer

	inventoryGRPC   *grpc.ClientConn
	inventoryClient client.InventoryClient

	paymentGRPC   *grpc.ClientConn
	paymentClient client.PaymentClient

	iamGRPCConn *grpc.ClientConn
	iamClient   auth_v1.AuthServiceClient

	pgxCon    *pgx.Conn
	db        *sql.DB
	orderRepo repoOrderInterface.OrderRepository

	orderService  serviceOrderInterface.OrderService
	orderV1Api    order_v1.Handler
	orderV1Server *order_v1.Server

	migrationRunner migrator.Migrator

	orderConsumerService      serviceOrderInterface.OrderConsumerService
	orderConsumer             wrappedKafka.Consumer
	consumerGroup             sarama.ConsumerGroup
	orderShipAssembledDecoder kafka.OrderAssembledDecoder
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (di *diContainer) IAMGRPCConn(_ context.Context) *grpc.ClientConn {
	if di.iamGRPCConn == nil {
		conn, err := grpc.NewClient(
			config.ApppConfig().IAMGRPCClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		if err != nil {
			panic(fmt.Sprintf("💥 failed to connect to payment service: %v", err))
		}

		closer.AddNamed("IAM gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		di.iamGRPCConn = conn
	}

	return di.iamGRPCConn
}

func (di *diContainer) IAMClient(ctx context.Context) auth_v1.AuthServiceClient {
	if di.iamClient == nil {
		di.iamClient = auth_v1.NewAuthServiceClient(di.IAMGRPCConn(ctx))
	}

	return di.iamClient
}

func (di *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if di.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.ApppConfig().Kafka.Brokers(),
			config.ApppConfig().OrderAssembledConsumerConfig.GroupID(),
			config.ApppConfig().OrderAssembledConsumerConfig.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return di.consumerGroup.Close()
		})
		di.consumerGroup = consumerGroup

	}
	return di.consumerGroup
}

func (di *diContainer) OrderConsumer(ctx context.Context) wrappedKafka.Consumer {
	if di.orderConsumer == nil {
		di.orderConsumer = wrappedKafkaConsumer.NewConsumer(
			di.ConsumerGroup(),
			[]string{
				config.ApppConfig().OrderAssembledConsumerConfig.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}
	return di.orderConsumer
}

func (di *diContainer) OrderConsumerService(ctx context.Context) serviceOrderInterface.OrderConsumerService {
	if di.orderConsumerService == nil {
		di.orderConsumerService = orderconsumer.NewService(
			di.OrderConsumer(ctx),
			di.OrderRepo(ctx),
			di.OrderShipAssembledDecoder(),
		)
	}
	return di.orderConsumerService
}

func (di *diContainer) OrderShipAssembledDecoder() kafka.OrderAssembledDecoder {
	if di.orderShipAssembledDecoder == nil {
		di.orderShipAssembledDecoder = decoder.NewOrderAssembledDecoder()
	}
	return di.orderShipAssembledDecoder
}

func (di *diContainer) MigrationRunner(ctx context.Context) migrator.Migrator {
	if di.migrationRunner == nil {
		di.migrationRunner = pgMigrator.NewMigrator(di.SqlDB(ctx), config.ApppConfig().Postges.MigrationDir())
	}
	return di.migrationRunner
}

func (di *diContainer) OrderV1Server(ctx context.Context) *order_v1.Server {
	if di.orderV1Server == nil {
		orderServer, err := order_v1.NewServer(di.OrderV1Api(ctx))
		if err != nil {
			panic(fmt.Sprintf("💥 failed to create OrderV1 server: %v", err))
		}
		di.orderV1Server = orderServer
	}
	return di.orderV1Server
}

func (di *diContainer) OrderV1Api(ctx context.Context) order_v1.Handler {
	if di.orderV1Api == nil {
		di.orderV1Api = serviceHandler.NewApi(di.OrderService(ctx))
	}
	return di.orderV1Api
}

func (di *diContainer) PaymentGRPC(ctx context.Context) *grpc.ClientConn {
	if di.paymentGRPC == nil {
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

		di.paymentGRPC = conn
	}

	return di.paymentGRPC
}

func (di *diContainer) PaymentClient(ctx context.Context) client.PaymentClient {
	if di.paymentClient == nil {
		di.paymentClient = paymentClient.NewClient(payment_v1.NewPaymentServiceClient(di.PaymentGRPC(ctx)))
	}

	return di.paymentClient
}

func (di *diContainer) InventoryGRPC(ctx context.Context) *grpc.ClientConn {
	if di.inventoryGRPC == nil {
		inventoryConn, err := grpc.NewClient(
			config.ApppConfig().InventoryGRPCClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("💥 failed to connect to inventory service: %v", err))
		}

		closer.AddNamed("Inventory gRPC connection", func(ctx context.Context) error {
			return inventoryConn.Close()
		})

		di.inventoryGRPC = inventoryConn
	}
	return di.inventoryGRPC
}

func (di *diContainer) InventoryClient(ctx context.Context) client.InventoryClient {
	if di.inventoryClient == nil {
		di.inventoryClient = inventoryClient.NewClient(inventory_v1.NewInventoryServiceClient(di.InventoryGRPC(ctx)))
	}
	return di.inventoryClient
}

func (di *diContainer) SyncProducer(ctx context.Context) sarama.SyncProducer {
	if di.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.ApppConfig().Kafka.Brokers(),
			config.ApppConfig().OrderPaidProducerConfig.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		di.syncProducer = p
	}
	return di.syncProducer
}

func (di *diContainer) OrderProducer(ctx context.Context) wrappedKafka.Producer {
	if di.orderProducer == nil {
		di.orderProducer = wrappedKafkaProducer.NewProducer(
			di.SyncProducer(ctx),
			config.ApppConfig().OrderPaidProducerConfig.Topic(),
			logger.Logger(),
		)
	}
	return di.orderProducer
}

func (di *diContainer) OrderProducerService(ctx context.Context) serviceOrderInterface.OrderProducerService {
	if di.orderProducerService == nil {
		di.orderProducerService = orderProducer.NewService(di.OrderProducer(ctx))
	}
	return di.orderProducerService
}

func (di *diContainer) OrderService(ctx context.Context) serviceOrderInterface.OrderService {
	if di.orderService == nil {
		di.orderService = serviceOrder.NewService(
			di.OrderRepo(ctx),
			di.PaymentClient(ctx),
			di.InventoryClient(ctx),
			di.OrderProducerService(ctx),
		)
	}
	return di.orderService
}

func (di *diContainer) PgxConnect(ctx context.Context) *pgx.Conn {
	if di.pgxCon == nil {
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

		di.pgxCon = conn
	}
	return di.pgxCon
}

func (di *diContainer) SqlDB(ctx context.Context) *sql.DB {
	if di.db == nil {
		di.db = stdlib.OpenDB(*di.PgxConnect(ctx).Config().Copy())
	}
	return di.db
}

func (di *diContainer) OrderRepo(ctx context.Context) repoOrderInterface.OrderRepository {
	if di.orderRepo == nil {
		di.orderRepo = repoOrder.NewRepository(di.SqlDB(ctx))
	}

	return di.orderRepo
}
