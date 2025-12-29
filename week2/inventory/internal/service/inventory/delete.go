package inventory

import "context"

func (s *InventoryService) DeleteByUUID(ctx context.Context, id string) error {
	if err := s.inventoryRepo.DeleteByUUID(ctx, id); err != nil {
		return err
	}
	return nil
}
