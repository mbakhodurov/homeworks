package inventory

import (
	"context"

	"github.com/mbakhodurov/homeworks/week2/inventory/internal/model"
	repoconverter "github.com/mbakhodurov/homeworks/week2/inventory/internal/repository/converter"
)

func (r *repository) GetAll(ctx context.Context) ([]model.Part, error) {

	r.mu.Lock()
	defer r.mu.Unlock()

	parts := make([]model.Part, 0, len(r.data))
	for _, part := range r.data {
		parts = append(parts, repoconverter.InventoryPartToModel(part))
	}
	return parts, nil
}
