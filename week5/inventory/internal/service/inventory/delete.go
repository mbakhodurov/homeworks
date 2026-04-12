package inventory

import (
	"context"

	"github.com/mbakhodurov/homeworks/week5/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) Delete(ctx context.Context, uuid string) (int64, error) {
	deletedCount, err := s.inventoryRepo.Delete(ctx, uuid)
	if err != nil {
		logger.Error(ctx, "failed to delete inventory",
			zap.String("uuid", uuid),
			zap.String("inventory-->Delete", ""),
			zap.Error(err),
		)
		return 0, err
	}

	return deletedCount, nil
}
