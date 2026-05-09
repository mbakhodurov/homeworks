package inventory

import (
	"context"

	"github.com/mbakhodurov/homeworks/week6/inventory/internal/model"
	"github.com/mbakhodurov/homeworks/week6/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) GetByUUID(ctx context.Context, uuid string) (model.Inventory, error) {
	res, err := s.inventoryRepo.GetByUUID(ctx, uuid)
	if err != nil {
		logger.Error(ctx, "failed to get inventory",
			zap.String("uuid", uuid),
			zap.String("inventory-->GetByUUID", ""),
			zap.Error(err),
		)
		return model.Inventory{}, err
	}
	return res, nil
}
