package v1

import (
	"context"
	"fmt"

	"github.com/mbakhodurov/homeworks/week3/inventory/internal/converter"
	inventory_v1 "github.com/mbakhodurov/homeworks/week3/shared/pkg/proto/inventory/v1"
)

func (a *InventoryApi) CreateInventory(ctx context.Context, req *inventory_v1.CreateInventoryRequest) (*inventory_v1.CreateInventoryResponse, error) {
	fmt.Println("req.GetInfo():", req.GetInfo())
	uuid, err := a.inventoryService.Create(ctx, converter.InventoryInfoProtoToModel(req.GetInfo()))
	if err != nil {
		return &inventory_v1.CreateInventoryResponse{}, err
	}
	return &inventory_v1.CreateInventoryResponse{
		Uuid: uuid,
	}, nil
}
