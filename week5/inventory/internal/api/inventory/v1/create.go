package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks/week5/inventory/internal/converter"
	"github.com/mbakhodurov/homeworks/week5/inventory/internal/model"
	inventory_v1 "github.com/mbakhodurov/homeworks/week5/shared/pkg/proto/inventory/v1"
)

func (a *InventoryV1Api) CreateInventory(ctx context.Context, req *inventory_v1.CreateInventoryRequest) (*inventory_v1.CreateInventoryResponse, error) {
	info := converter.InventoryInfoProtoToModel(req.GetInfo())

	inventory := model.Inventory{
		InventoryInfo: info,
	}
	uuid, err := a.inventoryService.Create(ctx, inventory)
	if err != nil {
		return &inventory_v1.CreateInventoryResponse{}, err
	}

	return &inventory_v1.CreateInventoryResponse{
		Uuid: uuid,
	}, nil
}
