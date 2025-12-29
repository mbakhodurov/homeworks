package service

import (
	"context"

	"github.com/mbakhodurov/homeworks/week2/inventory/internal/model"
)

type InventoryService interface {
	Create(ctx context.Context, part model.PartInfo) (string, error)
	GetAll(ctx context.Context) ([]model.Part, error)
	GetByUUID(ctx context.Context, id string) (model.Part, error)
	DeleteByUUID(ctx context.Context, id string) error
	Update(ctx context.Context, updatepart model.PartInfoUpdate, uuid string) error
	ListParts(ctx context.Context, parts model.PartsFilter) ([]model.Part, error)
}
