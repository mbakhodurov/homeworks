package order

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"

	client "github.com/mbakhodurov/homeworks/week7/order/internal/client/grpc"
	"github.com/mbakhodurov/homeworks/week7/order/internal/repository"
	def "github.com/mbakhodurov/homeworks/week7/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepo       repository.OrderRepository
	inventoryClient client.InventoryClient
	paymenClient    client.PaymentClient

	orderProducerService def.OrderProducerService

	ordersCounter  metric.Int64Counter
	revenueCounter metric.Float64Counter
}

func NewService(orderRepo repository.OrderRepository, paymentClient client.PaymentClient, inventoryClient client.InventoryClient, orderProducerService def.OrderProducerService) *service {
	meter := otel.Meter("order-service")

	ordersCounter, _ := meter.Int64Counter(
		"orders_total",
		metric.WithDescription("Total number of created orders"),
	)
	revenueCounter, _ := meter.Float64Counter(
		"orders_revenue_total",
		metric.WithDescription("Total revenue from created orders"),
	)

	return &service{
		orderRepo:            orderRepo,
		inventoryClient:      inventoryClient,
		paymenClient:         paymentClient,
		orderProducerService: orderProducerService,
		ordersCounter:        ordersCounter,
		revenueCounter:       revenueCounter,
	}
}
