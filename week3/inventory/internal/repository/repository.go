package repository

import (
	"context"

	"github.com/mbakhodurov/homeworks/week3/inventory/internal/model"
	"go.mongodb.org/mongo-driver/bson"
)

type InventoryRepository interface {
	Create(ctx context.Context, inventoryInfo model.Inventory) error
	GetAll(ctx context.Context) ([]model.Inventory, error)
	GetByUUID(ctx context.Context, uuid string) (model.Inventory, error)
	Delete(ctx context.Context, uuid string) (int64, error)
	Update(ctx context.Context, uuid string, set bson.M) error
	ListInventory(ctx context.Context, filter model.InventoryFilter) ([]model.Inventory, error)
}
