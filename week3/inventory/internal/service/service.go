package service

import (
	"context"

	"github.com/mbakhodurov/homeworks/week3/inventory/internal/model"
)

type InventoryService interface {
	Create(ctx context.Context, part model.InventoryInfo) (string, error)
	GetAll(ctx context.Context) ([]model.Inventory, error)
	GetByUUID(ctx context.Context, uuid string) (model.Inventory, error)
	Delete(ctx context.Context, uuid string) (int64, error)
	Update(ctx context.Context, uuid string, inventortUpdateInfo model.InventoryInfoUpdate) error
	ListParts(ctx context.Context, filter model.InventoryFilter) ([]model.Inventory, error)
}
