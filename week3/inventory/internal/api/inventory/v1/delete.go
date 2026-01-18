package v1

import (
	"context"

	inventory_v1 "github.com/mbakhodurov/homeworks/week3/shared/pkg/proto/inventory/v1"
)

func (a *InventoryApi) DeleteInventory(ctx context.Context, req *inventory_v1.DeleteInventoryRequest) (*inventory_v1.DeleteResponse, error) {
	deletedCount, err := a.inventoryService.Delete(ctx, req.GetUuid())
	if err != nil {
		return nil, err
	}

	return &inventory_v1.DeleteResponse{
		DeletedCount: deletedCount,
	}, nil
}
