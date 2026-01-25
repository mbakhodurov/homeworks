package inventory

import (
	"context"
	"errors"

	"github.com/mbakhodurov/homeworks/week3/inventory/internal/model"
	"github.com/mbakhodurov/homeworks/week3/inventory/internal/repository/converter"
	repoModel "github.com/mbakhodurov/homeworks/week3/inventory/internal/repository/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *repository) GetByUUID(ctx context.Context, uuid string) (model.Inventory, error) {

	var inventory repoModel.Inventory
	if err := r.collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&inventory); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Inventory{}, model.ErrInventoryNotFound
		}
		return model.Inventory{}, err
	}

	return converter.ConvertInventoryRepoModelToModel(inventory), nil
}
