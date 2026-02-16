package inventory

import (
	"context"

	"github.com/mbakhodurov/homeworks/week4/inventory/internal/model"
)

func (s *service) GetAll(ctx context.Context) ([]model.Inventory, error) {
	parts, err := s.inventoryRepo.GetAll(ctx)
	if err != nil {
		return []model.Inventory{}, err
	}
	return parts, nil
}
