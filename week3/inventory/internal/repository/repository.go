package repository

import (
	"context"

	"github.com/mbakhodurov/homeworks/week3/inventory/internal/model"
)

type InventoryRepository interface {
	Create(ctx context.Context, inventoryInfo model.InventoryInfo) (string, error)
	GetAll(ctx context.Context) ([]model.Inventory, error)
	GetByUUID(ctx context.Context, uuid string) (model.Inventory, error)
	Delete(ctx context.Context, uuid string) (int64, error)
	Update(ctx context.Context, uuid string, inventortUpdate model.PartInfoUpdate) error
}
