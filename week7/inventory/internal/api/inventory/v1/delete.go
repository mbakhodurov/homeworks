package v1

import (
	"context"

	inventory_v1 "github.com/mbakhodurov/homeworks/week7/shared/pkg/proto/inventory/v1"
)

func (a *InventoryV1Api) DeleteInventory(ctx context.Context, req *inventory_v1.DeleteInventoryRequest) (*inventory_v1.DeleteResponse, error) {
	deletedCount, err := a.inventoryService.Delete(ctx, req.GetUuid())
	if err != nil {
		return nil, err
	}

	return &inventory_v1.DeleteResponse{
		DeletedCount: deletedCount,
	}, nil
}
