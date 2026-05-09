package inventory

import (
	"context"

	"github.com/mbakhodurov/homeworks/week6/inventory/internal/model"
	"github.com/mbakhodurov/homeworks/week6/platform/pkg/logger"

	// grpcMid
	//
	// dleware "github.com/mbakhodurov/homeworks/week6/platform/pkg/middleware/grpc"
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
