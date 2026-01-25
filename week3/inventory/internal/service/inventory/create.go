package inventory

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mbakhodurov/homeworks/week3/inventory/internal/model"
)

func (s *service) Create(ctx context.Context, part model.InventoryInfo) (string, error) {
	// fmt.Println("part:", part)
	inventory := model.Inventory{
		UUID:          uuid.NewString(),
		InventoryInfo: part,
		CreatedAt:     time.Now(),
	}

	err := s.inventoryRepo.Create(ctx, inventory)
	if err != nil {
		return "", err
	}

	return inventory.UUID, err
}
