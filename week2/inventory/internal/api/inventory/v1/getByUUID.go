package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks/week2/inventory/internal/converter"
	inventory_v1 "github.com/mbakhodurov/homeworks/week2/shared/pkg/proto/inventory/v1"
)

func (a *InventoryApi) GetPart(ctx context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	inventory, err := a.inventoryService.GetByUUID(ctx, req.GetUuid())
	if err != nil {
		// if errors.Is(err, model.ErrInventoryNotFound) {
		// 	return nil, status.Errorf(codes.NotFound, "sighting with UUID %s not found", req.GetUuid())
		// }
		return nil, err
	}

	return &inventory_v1.GetPartResponse{
		Part: converter.InvertoryPartToProto(inventory),
	}, nil
}
