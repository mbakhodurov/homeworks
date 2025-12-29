package inventory

import (
	"context"

	"github.com/mbakhodurov/homeworks/week2/inventory/internal/model"
	"github.com/mbakhodurov/homeworks/week2/inventory/internal/repository/converter"
)

func (r *repository) GetByUUID(ctx context.Context, id string) (model.Part, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	inventory, ok := r.data[id]
	if !ok {
		return model.Part{}, model.ErrInventoryNotFound
	}
	return converter.InventoryPartToModel(inventory), nil
}
