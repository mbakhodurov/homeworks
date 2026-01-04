package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks/week2/order/internal/client/converter"
	"github.com/mbakhodurov/homeworks/week2/order/internal/model"
	inventory_v1 "github.com/mbakhodurov/homeworks/week2/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	parts, err := c.generatedClient.ListParts(ctx, &inventory_v1.ListPartsRequest{
		Filter: converter.PartsFilterToProto(filter),
	})

	if err != nil {
		return []model.Part{}, err
	}
	return converter.PartListToModel(parts), nil
}
