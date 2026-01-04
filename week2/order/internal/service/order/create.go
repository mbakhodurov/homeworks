package order

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mbakhodurov/homeworks/week2/order/internal/model"
)

const inventoryTimeout = 5 * time.Second

func (s *service) Create(ctx context.Context, user_uuid string, partUUIDs []string) (*model.Order, error) {
	if len(partUUIDs) == 0 {
		return nil, model.ErrPartsNotFound
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, inventoryTimeout)
	defer cancel()

	parts, err := s.inventoryClient.ListParts(ctxWithTimeout, model.PartsFilter{
		Uuids: &partUUIDs,
	})
	if err != nil {
		return nil, err
	}

	if len(parts) != len(partUUIDs) {
		return nil, model.ErrPartsNotFound
	}

	var totalPrice float32
	for _, part := range parts {
		totalPrice += float32(part.Partinfo.Price)
	}

	order := model.Order{
		Order_uuid:  uuid.NewString(),
		User_uuid:   user_uuid,
		Part_uuids:  partUUIDs,
		Total_price: float64(totalPrice),
		Status:      model.StatusPendingPayment,
		Created_at:  time.Now(),
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	return &order, nil
}
