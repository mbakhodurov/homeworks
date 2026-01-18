package v1

import (
	"github.com/mbakhodurov/homeworks/week3/inventory/internal/service"
	inventory_v1 "github.com/mbakhodurov/homeworks/week3/shared/pkg/proto/inventory/v1"
)

type InventoryApi struct {
	inventoryService service.InventoryService
	inventory_v1.UnimplementedInventoryServiceServer
}

func NewInventoryApi(inventoryService service.InventoryService) *InventoryApi {
	return &InventoryApi{
		inventoryService: inventoryService,
	}
}
