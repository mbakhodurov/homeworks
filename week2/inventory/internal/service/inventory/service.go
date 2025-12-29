package inventory

import (
	"github.com/mbakhodurov/homeworks/week2/inventory/internal/repository"
	def "github.com/mbakhodurov/homeworks/week2/inventory/internal/service"
)

var _ def.InventoryService = (*InventoryService)(nil)

type InventoryService struct {
	inventoryRepo repository.PartRepository
}

func NewInventoryServiceClient(inventoryRepo repository.PartRepository) *InventoryService {
	return &InventoryService{
		inventoryRepo: inventoryRepo,
	}
}
