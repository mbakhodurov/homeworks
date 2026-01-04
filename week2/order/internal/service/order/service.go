package order

import (
	"github.com/mbakhodurov/homeworks/week2/order/internal/client/grpc"
	"github.com/mbakhodurov/homeworks/week2/order/internal/repository"
	def "github.com/mbakhodurov/homeworks/week2/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient

	orderRepo repository.OrderRepository
}

func NewService(inventoryClient grpc.InventoryClient, paymentClient grpc.PaymentClient, orderRepo repository.OrderRepository) *service {
	return &service{
		orderRepo:       orderRepo,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
