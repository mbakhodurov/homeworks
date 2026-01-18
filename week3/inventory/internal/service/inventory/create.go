package inventory

import (
	"context"

	"github.com/mbakhodurov/homeworks/week3/inventory/internal/model"
)

func (s *InventoryService) Create(ctx context.Context, part model.InventoryInfo) (string, error) {
	uuid, err := s.inventoryRepo.Create(ctx, part)
	if err != nil {
		return "", err
	}
	return uuid, err
}
