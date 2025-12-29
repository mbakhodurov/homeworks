package inventory

import (
	"context"
	"time"

	"github.com/mbakhodurov/homeworks/week2/inventory/internal/model"
	repomodel "github.com/mbakhodurov/homeworks/week2/inventory/internal/repository/model"
	"github.com/samber/lo"
)

func (r *repository) Update(ctx context.Context, updateInfo model.PartInfoUpdate, uuid string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	inventory, ok := r.data[uuid]
	if !ok {
		return model.ErrInventoryNotFound
	}
	if updateInfo.Name != nil {
		inventory.Partinfo.Name = *updateInfo.Name
	}
	if updateInfo.Description != nil {
		inventory.Partinfo.Description = *updateInfo.Description
	}
	if updateInfo.Price != nil {
		inventory.Partinfo.Price = *updateInfo.Price
	}
	if updateInfo.StockQuantity != nil {
		inventory.Partinfo.Stock_quantity = *updateInfo.StockQuantity
	}

	if updateInfo.Category != nil {
		inventory.Partinfo.Category = repomodel.Category(*updateInfo.Category)
	}
	if updateInfo.Dimensions != nil {
		inventory.Partinfo.Dimensions = repomodel.Dimensions(*updateInfo.Dimensions)
	}
	if updateInfo.Manufacturer != nil {
		inventory.Partinfo.Manufacturer = repomodel.Manufacturer(*updateInfo.Manufacturer)
	}
	inventory.UpdatedAt = lo.ToPtr(time.Now())
	r.data[uuid] = inventory
	return nil
}
