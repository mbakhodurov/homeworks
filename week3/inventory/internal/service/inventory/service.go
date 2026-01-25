package inventory

import (
	"github.com/mbakhodurov/homeworks/week3/inventory/internal/repository"
	def "github.com/mbakhodurov/homeworks/week3/inventory/internal/service"
)

var _ def.InventoryService = (*service)(nil)

type service struct {
	inventoryRepo repository.InventoryRepository
}

func NewService(inventoryRepo repository.InventoryRepository) *service {
	return &service{
		inventoryRepo: inventoryRepo,
	}
}
