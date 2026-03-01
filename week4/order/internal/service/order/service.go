package order

import (
	"github.com/mbakhodurov/homeworks/week4/order/internal/client/grpc"
	"github.com/mbakhodurov/homeworks/week4/order/internal/repository"
	def "github.com/mbakhodurov/homeworks/week4/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepo       repository.OrderRepository
	paymentClient   grpc.PaymentClient
	inventoryClient grpc.InventoryClient
}

func NewService(orderRepo repository.OrderRepository, inventoryClient grpc.InventoryClient, paymentClient grpc.PaymentClient) *service {
	return &service{
		orderRepo:       orderRepo,
		paymentClient:   paymentClient,
		inventoryClient: inventoryClient,
	}
}
