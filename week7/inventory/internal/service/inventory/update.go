package inventory

import (
	"context"
	"errors"
	"time"

	"github.com/mbakhodurov/homeworks/week7/inventory/internal/model"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func (s *service) Update(ctx context.Context, uuid string, inventortUpdateInfo model.InventoryInfoUpdate) error {
	set := bson.M{}
	_, err := s.inventoryRepo.GetByUUID(ctx, uuid)
	if err != nil {
		if errors.Is(err, model.ErrInventoryNotFound) {
			logger.Error(ctx, "failed to Update inventory",
				zap.Any("inventortUpdateInfo", inventortUpdateInfo),
				zap.String("inventory-->inventoryRepo.GetByUUID", ""),
				zap.Error(err),
			)
			return model.ErrInventoryNotFound
		}
		return err
	}

	if inventortUpdateInfo.Name != nil {
		set["inventory_info.name"] = *inventortUpdateInfo.Name
	}
	if inventortUpdateInfo.Description != nil {
		set["inventory_info.description"] = *inventortUpdateInfo.Description
	}
	if inventortUpdateInfo.Price != nil {
		set["inventory_info.price"] = *inventortUpdateInfo.Price
	}
	if inventortUpdateInfo.StockQuantity != nil {
		set["inventory_info.stock_quantity"] = *inventortUpdateInfo.StockQuantity
	}
	if inventortUpdateInfo.Category != nil {
		set["inventory_info.category"] = *inventortUpdateInfo.Category
	}
	if inventortUpdateInfo.Dimensions != nil {
		set["inventory_info.dimensions"] = *inventortUpdateInfo.Dimensions
	}
	if inventortUpdateInfo.Manufacturer != nil {
		set["inventory_info.manufacturer"] = *inventortUpdateInfo.Manufacturer
	}

	if len(set) == 0 {
		return nil
	}

	set["updated_at"] = time.Now()

	err = s.inventoryRepo.Update(ctx, uuid, set)
	if err != nil {
		logger.Error(ctx, "failed to Update inventory",
			zap.Any("inventortUpdateInfo", inventortUpdateInfo),
			zap.String("inventory-->inventortUpdateInfo", ""),
			zap.Error(err),
		)
		return err
	}

	return nil
}
