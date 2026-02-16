package inventory

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mbakhodurov/homeworks/week4/inventory/internal/model"
)

func (s *service) Create(ctx context.Context, info model.Inventory) (string, error) {
	info.UUID = uuid.NewString()
	info.CreatedAt = time.Now()

	err := s.inventoryRepo.Create(ctx, info)
	if err != nil {
		return "", err
	}
	return info.UUID, nil
}
