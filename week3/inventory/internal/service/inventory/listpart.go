package inventory

import (
	"context"

	"github.com/mbakhodurov/homeworks/week3/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filter model.InventoryFilter) ([]model.Inventory, error) {
	parts, err := s.inventoryRepo.ListInventory(ctx, filter)
	if err != nil {
		return nil, err
	}
	return parts, nil
}
