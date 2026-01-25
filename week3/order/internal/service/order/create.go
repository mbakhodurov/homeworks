package order

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mbakhodurov/homeworks/week3/order/internal/model"
	"github.com/samber/lo"
)

const inventoryTimeout = 5 * time.Second

func (s *service) Create(ctx context.Context, user_uuid string, partUUIDs []string) (*model.Order, error) {
	if len(partUUIDs) == 0 {
		return nil, model.ErrPartsNotFound
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, inventoryTimeout)
	defer cancel()

	inventoryParts, err := s.inventoryClient.InventoryPart(ctxWithTimeout, model.InventoryFilter{
		UUID: &partUUIDs,
	})

	if err != nil {
		return nil, err
	}

	if len(inventoryParts) != len(partUUIDs) {
		return nil, model.ErrPartsNotFound
	}

	var totalPrice float32
	for _, inventoryPart := range inventoryParts {
		totalPrice += float32(inventoryPart.InventoryInfo.Price)
	}

	order := model.Order{
		OrderUUID:  uuid.NewString(),
		UserUUID:   user_uuid,
		PartUUIDs:  partUUIDs,
		TotalPrice: float64(totalPrice),
		Status:     model.StatusPendingPayment,
		CreatedAt:  lo.ToPtr(time.Now()),
		Updated_at: lo.ToPtr(time.Now()),
	}

	_, err = s.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}
