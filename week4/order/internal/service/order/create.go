package order

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mbakhodurov/homeworks/week4/order/internal/model"
	"github.com/mbakhodurov/homeworks/week4/platform/pkg/logger"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

const inventoryTimeout = 5 * time.Second

func (s *service) Create(ctx context.Context, user_uuid string, partUUIDs []string) (*model.Order, error) {
	if len(partUUIDs) == 0 {
		logger.Error(ctx, "PartUUID отсутствует")
		return nil, model.ErrPartsNotFound
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, inventoryTimeout)
	defer cancel()

	inventories, err := s.inventoryClient.ListPartInventory(ctxWithTimeout, model.InventoryFilter{
		UUID: &partUUIDs,
	})
	if err != nil {
		logger.Error(ctx, "Ошибка с работой с InventoryClient", zap.Error(err))
		return &model.Order{}, err
	}

	if len(inventories) != len(partUUIDs) {
		logger.Error(ctx, "Ошибка с работой с InventoryClient", zap.Error(err))
		return nil, model.ErrPartsNotFound
	}

	var sum float64
	for _, v := range inventories {
		sum += v.InventoryInfo.Price
	}

	order := model.Order{
		OrderUUID:  uuid.NewString(),
		UserUUID:   user_uuid,
		PartUUID:   partUUIDs,
		TotalPrice: sum,
		Status:     model.StatusPendingPayment,
		CreatedAt:  lo.ToPtr(time.Now()),
		UpdatedAt:  lo.ToPtr(time.Now()),
	}

	_, err = s.orderRepo.CreateOrder(ctx, order)
	if err != nil {
		logger.Error(ctx, "Ошибка с работой с CreateOrder", zap.Error(err))
		return &model.Order{}, err
	}

	return &order, nil
}
