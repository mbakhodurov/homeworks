package inventory

import (
	"github.com/mbakhodurov/homeworks/week3/inventory/internal/repository"
	def "github.com/mbakhodurov/homeworks/week3/inventory/internal/service"
)

var _ def.InventoryService = (*InventoryService)(nil)

type InventoryService struct {
	inventoryRepo repository.InventoryRepository
}

func NewInventoryService(inventoryRepo repository.InventoryRepository) *InventoryService {
	return &InventoryService{
		inventoryRepo: inventoryRepo,
	}
}
