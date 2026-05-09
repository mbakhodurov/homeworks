package v1

import (
	"github.com/mbakhodurov/homeworks/week6/inventory/internal/service"
	inventory_v1 "github.com/mbakhodurov/homeworks/week6/shared/pkg/proto/inventory/v1"
)

type InventoryV1Api struct {
	inventory_v1.UnimplementedInventoryServiceServer
	inventoryService service.InventoryService
}

func NewInventoryV1Api(inventoryService service.InventoryService) *InventoryV1Api {
	return &InventoryV1Api{
		inventoryService: inventoryService,
	}
}
