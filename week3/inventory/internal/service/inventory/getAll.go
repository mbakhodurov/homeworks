package inventory

import (
	"context"

	"github.com/mbakhodurov/homeworks/week3/inventory/internal/model"
)

func (s *InventoryService) GetAll(ctx context.Context) ([]model.Inventory, error) {
	parts, err := s.inventoryRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return parts, nil
}
