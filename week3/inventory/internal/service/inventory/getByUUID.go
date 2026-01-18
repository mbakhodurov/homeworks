package inventory

import (
	"context"

	"github.com/mbakhodurov/homeworks/week3/inventory/internal/model"
)

func (s *InventoryService) GetByUUID(ctx context.Context, uuid string) (model.Inventory, error) {
	inventory, err := s.inventoryRepo.GetByUUID(ctx, uuid)
	if err != nil {
		return model.Inventory{}, err
	}

	return inventory, nil
}
