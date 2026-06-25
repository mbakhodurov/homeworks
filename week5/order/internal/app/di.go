package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	client "github.com/mbakhodurov/homeworks/week5/order/internal/client/grpc"
	inventoryClient "github.com/mbakhodurov/homeworks/week5/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/mbakhodurov/homeworks/week5/order/internal/client/grpc/payment/v1"
	"github.com/mbakhodurov/homeworks/week5/order/internal/converter/kafka"
	"github.com/mbakhodurov/homeworks/week5/order/internal/converter/kafka/decoder"
	repoOrderInterface "github.com/mbakhodurov/homeworks/week5/order/internal/repository"
	repo "github.com/mbakhodurov/homeworks/week5/order/internal/repository/order"
	"github.com/mbakhodurov/homeworks/week5/order/internal/service"

	serviceHandler "github.com/mbakhodurov/homeworks/week5/order/internal/api/order/v1"
	serviceOrderInterface "github.com/mbakhodurov/homeworks/week5/order/internal/service"
	orderconsumer "github.com/mbakhodurov/homeworks/week5/order/internal/service/consumer/order_consumer"
	serviceOrder "github.com/mbakhodurov/homeworks/week5/order/internal/service/order"

	"github.com/mbakhodurov/homeworks/week5/order/internal/config"
	orderProducer "github.com/mbakhodurov/homeworks/week5/order/internal/service/producer/order_producer"
	"github.com/mbakhodurov/homeworks/week5/platform/pkg/closer"
	wrappedKafka "github.com/mbakhodurov/homeworks/week5/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/mbakhodurov/homeworks/week5/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/mbakhodurov/homeworks/week5/platform/pkg/kafka/producer"

	"github.com/mbakhodurov/homeworks/week5/platform/pkg/logger"

	"github.com/mbakhodurov/homeworks/week5/platform/pkg/migrator"

	pgMigrator "github.com/mbakhodurov/homeworks/week5/platform/pkg/migrator/pg"

	kafkaMiddleware "github.com/mbakhodurov/homeworks/week5/platform/pkg/middleware/kafka"
	order_v1 "github.com/mbakhodurov/homeworks/week5/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/mbakhodurov/homeworks/week5/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/mbakhodurov/homeworks/week5/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type diContainer struct {
	inventoryGRPC   *grpc.ClientConn
	inventoryClient client.InventoryClient

	paymentGRPC   *grpc.ClientConn
	paymentClient client.PaymentClient

	pgxCon          *pgx.Conn
	sqlDb           *sql.DB
	orderRepository repoOrderInterface.OrderRepository

	syncProducer         sarama.SyncProducer
	orderProducer        wrappedKafka.Producer
	orderProducerService service.OrderProducerService

	orderService serviceOrderInterface.OrderService

	migrationRunner migrator.Migrator

	orderV1API    order_v1.Handler
	orderV1Server *order_v1.Server

	config config.LoggerConfig

	orderConsumerService      service.OrderConsumerService
	orderConsumer             wrappedKafka.Consumer
	consumerGroup             sarama.ConsumerGroup
	orderShipAssembledDecoder kafka.OrderAssembledDecoder
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (di *diContainer) OrderShipAssembledDecoder() kafka.OrderAssembledDecoder {
	if di.orderShipAssembledDecoder == nil {
		di.orderShipAssembledDecoder = decoder.NewOrderAssembledDecoder()
	}
	return di.orderShipAssembledDecoder
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

func (di *diContainer) OrderConsumer() wrappedKafka.Consumer {
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

func (di *diContainer) OrderConsumerService(ctx context.Context) service.OrderConsumerService {
	if di.orderConsumerService == nil {
		di.orderConsumerService = orderconsumer.NewService(
			di.OrderConsumer(),
			di.OrderRepository(ctx),
			di.OrderShipAssembledDecoder(),
		)
	}
	return di.orderConsumerService
}

func (di *diContainer) MigrationRunner(ctx context.Context) migrator.Migrator {
	if di.migrationRunner == nil {
		di.migrationRunner = pgMigrator.NewMigrator(di.SqlDB(ctx), config.ApppConfig().Postges.MigrationDir())
	}
	return di.migrationRunner
}

func (di *diContainer) OrderV1Server(ctx context.Context) *order_v1.Server {
	if di.orderV1Server == nil {
		orderServer, err := order_v1.NewServer(di.OrderV1API(ctx))
		if err != nil {
			panic(fmt.Sprintf("💥 failed to create OrderV1 server: %v", err))
		}
		di.orderV1Server = orderServer
	}
	return di.orderV1Server
}

func (di *diContainer) OrderV1API(ctx context.Context) order_v1.Handler {
	if di.orderV1API == nil {
		di.orderV1API = serviceHandler.NewApi(di.OrderService(ctx))
	}

	return di.orderV1API
}

func (di *diContainer) OrderProducerService(ctx context.Context) service.OrderProducerService {
	if di.orderProducerService == nil {
		di.orderProducerService = orderProducer.NewService(di.OrderProducer(ctx))
	}
	return di.orderProducerService
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

func (di *diContainer) SyncProducer(ctx context.Context) sarama.SyncProducer {
	if di.syncProducer == nil {
		// configs := sarama.NewConfig()
		// configs.Producer.RequiredAcks = sarama.WaitForAll
		// configs.Producer.Retry.Max = 5
		// configs.Producer.Return.Successes = true

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

func (di *diContainer) OrderService(ctx context.Context) serviceOrderInterface.OrderService {
	if di.orderService == nil {
		di.orderService = serviceOrder.NewService(di.OrderRepository(ctx), di.PaymentClient(ctx), di.InventoryClient(ctx), di.OrderProducerService(ctx))
	}
	return di.orderService
}

func (di *diContainer) OrderRepository(ctx context.Context) repoOrderInterface.OrderRepository {
	if di.orderRepository == nil {
		di.orderRepository = repo.NewRepository(di.SqlDB(ctx))
	}
	return di.orderRepository
}

func (di *diContainer) SqlDB(ctx context.Context) *sql.DB {
	if di.sqlDb == nil {
		di.sqlDb = stdlib.OpenDB(*di.PgxConnect(ctx).Config().Copy())
	}
	return di.sqlDb
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
