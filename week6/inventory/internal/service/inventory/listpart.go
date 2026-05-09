package inventory

import (
	"context"

	"github.com/mbakhodurov/homeworks/week6/inventory/internal/model"
	"github.com/mbakhodurov/homeworks/week6/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) ListParts(ctx context.Context, filter model.InventoryFilter) ([]model.Inventory, error) {
	parts, err := s.inventoryRepo.ListParts(ctx, filter)
	if err != nil {
		logger.Error(ctx, "failed to ListParts inventory",
			zap.String("inventory-->ListParts", ""),
			zap.Error(err),
		)
		return nil, err
	}
	return parts, nil
}
