package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks/week2/inventory/internal/converter"
	inventory_v1 "github.com/mbakhodurov/homeworks/week2/shared/pkg/proto/inventory/v1"
)

func (a *InventoryApi) GetAll(ctx context.Context, req *inventory_v1.GetAllRequest) (*inventory_v1.GetAllResponse, error) {
	sightings, err := a.inventoryService.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	protoinventories := make([]*inventory_v1.Part, 0, len(sightings))
	for _, v := range sightings {
		protoinventories = append(protoinventories, converter.InvertoryPartToProto(v))
	}

	return &inventory_v1.GetAllResponse{
		Part:       protoinventories,
		TotalCount: int64(len(protoinventories)),
	}, nil
}
