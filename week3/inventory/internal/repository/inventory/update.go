package inventory

import (
	"context"
	"errors"
	"time"

	"github.com/mbakhodurov/homeworks/week3/inventory/internal/model"
	"github.com/mbakhodurov/homeworks/week3/inventory/internal/repository/converter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *repository) Update(ctx context.Context, uuid string, inventortUpdateInfo model.PartInfoUpdate) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	err := r.collection.FindOne(ctx, bson.M{"uuid": uuid}).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.ErrInventoryNotFound
		}
		return err
	}

	repoUpdateInfo := converter.ConvertInventoryUpdateModelToRepoModel(inventortUpdateInfo)

	set := bson.M{}

	if repoUpdateInfo.Name != nil {
		set["inventory_info.name"] = *repoUpdateInfo.Name
	}
	if repoUpdateInfo.Description != nil {
		set["inventory_info.description"] = *repoUpdateInfo.Description
	}
	if repoUpdateInfo.Price != nil {
		set["inventory_info.price"] = *repoUpdateInfo.Price
	}
	if repoUpdateInfo.StockQuantity != nil {
		set["inventory_info.stock_quantity"] = *repoUpdateInfo.StockQuantity
	}
	if repoUpdateInfo.Category != nil {
		set["inventory_info.category"] = *repoUpdateInfo.Category
	}
	if repoUpdateInfo.Dimensions != nil {
		set["inventory_info.dimensions"] = *repoUpdateInfo.Dimensions
	}
	if repoUpdateInfo.Manufacturer != nil {
		set["inventory_info.manufacturer"] = *repoUpdateInfo.Manufacturer
	}

	set["updated_at"] = time.Now()

	if len(set) == 0 {
		return nil
	}

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"uuid": uuid},
		bson.M{"$set": set},
	)
	if err != nil {
		return err
	}

	return nil
}
