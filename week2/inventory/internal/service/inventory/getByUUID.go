package inventory

import (
	"context"

	"github.com/mbakhodurov/homeworks/week2/inventory/internal/model"
)

func (s *InventoryService) GetByUUID(ctx context.Context, id string) (model.Part, error) {
	inventory, err := s.inventoryRepo.GetByUUID(ctx, id)
	if err != nil {
		return model.Part{}, err
	}
	return inventory, nil
}
