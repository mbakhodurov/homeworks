package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks/week3/inventory/internal/converter"
	inventory_v1 "github.com/mbakhodurov/homeworks/week3/shared/pkg/proto/inventory/v1"
)

func (a *InventoryApi) GetInventoryByUUID(ctx context.Context, req *inventory_v1.GetInventoryByUUIDRequest) (*inventory_v1.GetInventoryResponse, error) {
	inventory, err := a.inventoryService.GetByUUID(ctx, req.GetUuid())
	if err != nil {
		return nil, err
	}

	return &inventory_v1.GetInventoryResponse{
		Inventory: converter.InventoryModelToProto(inventory),
	}, nil
}
