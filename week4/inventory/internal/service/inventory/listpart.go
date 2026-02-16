package inventory

import (
	"context"

	"github.com/mbakhodurov/homeworks/week4/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filter model.InventoryFilter) ([]model.Inventory, error) {
	parts, err := s.inventoryRepo.ListParts(ctx, filter)
	if err != nil {
		return nil, err
	}
	return parts, nil
}
