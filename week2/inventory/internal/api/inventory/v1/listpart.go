package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks/week2/inventory/internal/converter"
	inventory_v1 "github.com/mbakhodurov/homeworks/week2/shared/pkg/proto/inventory/v1"
)

func (a *InventoryApi) ListParts(ctx context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	parts, err := a.inventoryService.ListParts(ctx, converter.ProtoToPartsFilter(req.GetFilter()))
	if err != nil {
		return nil, err
	}

	protoparts := make([]*inventory_v1.Part, 0, len(parts))
	for _, part := range parts {
		protoparts = append(protoparts, converter.InvertoryPartToProto(part))
	}

	return &inventory_v1.ListPartsResponse{
		Part:       protoparts,
		TotalCount: int64(len(protoparts)),
	}, nil
}
