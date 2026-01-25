package order

import (
	"github.com/mbakhodurov/homeworks/week3/order/internal/client/grpc"
	"github.com/mbakhodurov/homeworks/week3/order/internal/repository"
	def "github.com/mbakhodurov/homeworks/week3/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepo       repository.OrderRepository
	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient
}

func NewService(orderRepo repository.OrderRepository, inventoryClient grpc.InventoryClient, paymentClient grpc.PaymentClient) *service {
	return &service{
		orderRepo:       orderRepo,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
