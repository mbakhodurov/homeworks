package order

import (
	"context"

	"github.com/mbakhodurov/homeworks/week2/order/internal/model"
	repoModel "github.com/mbakhodurov/homeworks/week2/order/internal/repository/model"
)

func (r *repository) Update(ctx context.Context, uuid string, info model.OrderUpdateInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.data[uuid]
	if !ok {
		return model.ErrOrderNotFound
	}

	if info.Status != nil {
		order.Status = repoModel.OrderStatus(*info.Status)
	}

	if info.Payment_method != nil {
		pm := repoModel.PaymentMethod(*info.Payment_method)
		order.Payment_method = &pm
	}

	if info.Updated_at != nil {
		order.Updated_at = *info.Updated_at
	}

	if info.Deleted_at != nil {
		order.Deleted_at = *info.Deleted_at
	}

	r.data[uuid] = order

	return nil
}
