package v1

import (
	inventoryService "github.com/mbakhodurov/homeworks/week2/inventory/internal/service"
	inventory_v1 "github.com/mbakhodurov/homeworks/week2/shared/pkg/proto/inventory/v1"
)

type InventoryApi struct {
	inventoryService inventoryService.InventoryService
	inventory_v1.UnimplementedInventoryServiceServer
}

func NewInventoryApi(inventoryService inventoryService.InventoryService) *InventoryApi {
	return &InventoryApi{
		inventoryService: inventoryService,
	}
}
