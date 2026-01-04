package order

import (
	"context"

	"github.com/mbakhodurov/homeworks/week2/order/internal/model"
	"github.com/mbakhodurov/homeworks/week2/order/internal/repository/converter"
)

func (r *repository) Get(ctx context.Context, OrderUUID string) (model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.data[OrderUUID]
	if !ok {
		return model.Order{}, model.ErrOrderNotFound
	}

	return converter.OrderToModel(order), nil
}
