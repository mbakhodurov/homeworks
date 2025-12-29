package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks/week2/inventory/internal/converter"
	inventory_v1 "github.com/mbakhodurov/homeworks/week2/shared/pkg/proto/inventory/v1"
)

func (a *InventoryApi) Create(ctx context.Context, req *inventory_v1.CreatePartsRequest) (*inventory_v1.CreatePartsResponse, error) {
	uuid, err := a.inventoryService.Create(ctx, converter.InventoryPartInfoProtoToModel(req.GetInfo()))
	if err != nil {
		return nil, err
	}

	return &inventory_v1.CreatePartsResponse{
		Uuid: uuid,
	}, nil
}
