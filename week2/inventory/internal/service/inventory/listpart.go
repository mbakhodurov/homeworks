package inventory

import (
	"context"

	"github.com/mbakhodurov/homeworks/week2/inventory/internal/model"
)

func (s *InventoryService) ListParts(ctx context.Context, parts model.PartsFilter) ([]model.Part, error) {
	inventories, err := s.inventoryRepo.ListParts(ctx, parts)
	if err != nil {
		return nil, err
	}
	return inventories, nil
}
