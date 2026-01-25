package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks/week3/inventory/internal/converter"
	inventory_v1 "github.com/mbakhodurov/homeworks/week3/shared/pkg/proto/inventory/v1"
)

func (a *InventoryApi) ListPartInventory(ctx context.Context, req *inventory_v1.ListPartInventoryRequest) (*inventory_v1.ListPartInventoryResponse, error) {
	inventories, err := a.inventoryService.ListParts(ctx, converter.InventortListPartProtoToModel(req.GetFilter()))
	if err != nil {
		return nil, err
	}

	protoinventories := make([]*inventory_v1.Inventory, 0, len(inventories))
	for _, v := range inventories {
		protoinventories = append(protoinventories, converter.InventoryModelToProto(v))
	}
	return &inventory_v1.ListPartInventoryResponse{
		Inventory:  protoinventories,
		TotalCount: int64(len(protoinventories)),
	}, nil
}
