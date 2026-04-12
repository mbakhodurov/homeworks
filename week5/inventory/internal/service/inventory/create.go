package inventory

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mbakhodurov/homeworks/week5/inventory/internal/model"
	"github.com/mbakhodurov/homeworks/week5/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) Create(ctx context.Context, info model.Inventory) (string, error) {
	info.UUID = uuid.NewString()
	info.CreatedAt = time.Now()
	if err := s.inventoryRepo.Create(ctx, info); err != nil {
		logger.Error(ctx, "failed to create inventory",
			zap.String("uuid", info.UUID),
			zap.String("inventory-->Create", ""),
			zap.Error(err),
		)
		return "", err
	}

	return info.UUID, nil
}
