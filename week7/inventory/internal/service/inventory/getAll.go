package inventory

import (
	"context"

	"github.com/mbakhodurov/homeworks/week7/inventory/internal/model"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/logger"

	"go.uber.org/zap"
)

func (s *service) GetAll(ctx context.Context) ([]model.Inventory, error) {

	parts, err := s.inventoryRepo.GetAll(ctx)
	if err != nil {
		logger.Error(ctx, "failed to GetAll inventory",
			zap.String("inventory-->GetAll", ""),
			zap.Error(err),
		)
		return []model.Inventory{}, err
	}
	return parts, nil
}
