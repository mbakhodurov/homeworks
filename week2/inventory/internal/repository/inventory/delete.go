package inventory

import (
	"context"
	"time"

	"github.com/mbakhodurov/homeworks/week2/inventory/internal/model"
	"github.com/samber/lo"
)

func (r *repository) DeleteByUUID(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	inventory, ok := r.data[id]
	if !ok {
		return model.ErrInventoryNotFound
	}

	inventory.DeletedAt = lo.ToPtr(time.Now())
	r.data[id] = inventory
	return nil
}
