package order

import (
	"context"

	"github.com/mbakhodurov/homeworks/week2/order/internal/model"
	repoContever "github.com/mbakhodurov/homeworks/week2/order/internal/repository/converter"
)

func (r *repository) GetAll(ctx context.Context) ([]model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	orders := make([]model.Order, 0, len(r.data))

	for _, order := range r.data {
		orders = append(orders, repoContever.OrderToModel(order))
	}
	return orders, nil
}
