package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks/week3/inventory/internal/converter"
	inventory_v1 "github.com/mbakhodurov/homeworks/week3/shared/pkg/proto/inventory/v1"
)

func (a *InventoryApi) GetAllInventory(ctx context.Context, req *inventory_v1.GetAllInventoryRequest) (*inventory_v1.GetAllInventoryResponse, error) {
	inventories, err := a.inventoryService.GetAll(ctx)
	if err != nil {
		return &inventory_v1.GetAllInventoryResponse{}, err
	}

	protoinventories := make([]*inventory_v1.Inventory, 0, len(inventories))
	for _, v := range inventories {
		protoinventories = append(protoinventories, converter.InventoryModelToProto(v))
	}
	return &inventory_v1.GetAllInventoryResponse{
		Inventory:  protoinventories,
		TotalCount: int64(len(protoinventories)),
	}, nil
}
