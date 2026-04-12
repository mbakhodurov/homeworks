package order

import (
	client "github.com/mbakhodurov/homeworks/week5/order/internal/client/grpc"
	"github.com/mbakhodurov/homeworks/week5/order/internal/repository"
	def "github.com/mbakhodurov/homeworks/week5/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepo            repository.OrderRepository
	orderProducerService def.OrderProducerService

	paymentClient   client.PaymentClient
	inventoryClient client.InventoryClient
}

func NewService(orderRepo repository.OrderRepository, paymentClient client.PaymentClient, inventoryClient client.InventoryClient, orderProduceService def.OrderProducerService) *service {
	return &service{
		orderRepo:            orderRepo,
		paymentClient:        paymentClient,
		inventoryClient:      inventoryClient,
		orderProducerService: orderProduceService,
	}
}
