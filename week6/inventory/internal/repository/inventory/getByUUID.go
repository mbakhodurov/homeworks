package inventory

import (
	"context"
	"errors"

	"github.com/mbakhodurov/homeworks/week6/inventory/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *repository) GetByUUID(ctx context.Context, uuid string) (model.Inventory, error) {
	var inventory model.Inventory

	if err := r.collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(inventory); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Inventory{}, model.ErrInventoryNotFound
		}
		return model.Inventory{}, err
	}

	return inventory, nil
}
