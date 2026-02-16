package inventory

import (
	"context"
	"fmt"

	"github.com/mbakhodurov/homeworks/week4/inventory/internal/model"
	"github.com/mbakhodurov/homeworks/week4/inventory/internal/repository/converter"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *repository) Create(ctx context.Context, info model.Inventory) error {
	_, err := r.collection.InsertOne(ctx, converter.InventoryModelToRepoModel(info))
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("inventory already exists: %v", err)
		}
		return err
	}
	return nil
}
