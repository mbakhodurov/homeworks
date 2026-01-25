package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks/week3/order/internal/client/converter"
	"github.com/mbakhodurov/homeworks/week3/order/internal/model"
	inventory_v1 "github.com/mbakhodurov/homeworks/week3/shared/pkg/proto/inventory/v1"
)

func (c *client) InventoryPart(ctx context.Context, filter model.InventoryFilter) ([]model.Inventory, error) {
	res, err := c.generatedClient.ListPartInventory(ctx, &inventory_v1.ListPartInventoryRequest{
		Filter: converter.ModelFilterToProtoFilter(filter),
	})

	if err != nil {
		return []model.Inventory{}, err
	}

	return converter.ResponseProtoToInventoryModel(res), nil
}
