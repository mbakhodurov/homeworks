package inventory

import (
	"context"

	"github.com/mbakhodurov/homeworks/week2/inventory/internal/model"
)

func (s *InventoryService) GetAll(ctx context.Context) ([]model.Part, error) {
	inventories, err := s.inventoryRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return inventories, nil
}
