package order

import (
	client "github.com/mbakhodurov/homeworks/week6/order/internal/client/grpc"
	"github.com/mbakhodurov/homeworks/week6/order/internal/repository"
	def "github.com/mbakhodurov/homeworks/week6/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepo       repository.OrderRepository
	inventoryClient client.InventoryClient
	paymenClient    client.PaymentClient

	orderProducerService def.OrderProducerService
}

func NewService(orderRepo repository.OrderRepository, inventoryClient client.InventoryClient, paymenClient client.PaymentClient, orderProducerService def.OrderProducerService) *service {
	return &service{
		orderRepo:            orderRepo,
		inventoryClient:      inventoryClient,
		paymenClient:         paymenClient,
		orderProducerService: orderProducerService,
	}
}
