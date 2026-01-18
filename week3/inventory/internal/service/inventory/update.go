package inventory

import (
	"context"

	"github.com/mbakhodurov/homeworks/week3/inventory/internal/model"
)

func (s *InventoryService) Update(ctx context.Context, uuid string, inventortUpdateInfo model.PartInfoUpdate) error {
	err := s.inventoryRepo.Update(ctx, uuid, inventortUpdateInfo)
	if err != nil {
		return err
	}

	return nil
}
