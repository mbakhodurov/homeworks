package inventory

import "context"

func (s *service) Delete(ctx context.Context, uuid string) (int64, error) {
	deletedCount, err := s.inventoryRepo.Delete(ctx, uuid)
	if err != nil {
		return 0, err
	}

	return deletedCount, nil
}
