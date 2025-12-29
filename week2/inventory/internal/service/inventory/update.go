package inventory

import (
	"context"

	"github.com/mbakhodurov/homeworks/week2/inventory/internal/model"
)

func (s *InventoryService) Update(ctx context.Context, updatepart model.PartInfoUpdate, uuid string) error {
	if err := s.inventoryRepo.Update(ctx, updatepart, uuid); err != nil {
		return err
	}
	return nil
}
