package inventory

import (
	"context"
	"errors"
	"log"

	"github.com/mbakhodurov/homeworks/week3/inventory/internal/model"
	"github.com/mbakhodurov/homeworks/week3/inventory/internal/repository/converter"
	repomodel "github.com/mbakhodurov/homeworks/week3/inventory/internal/repository/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *repository) ListInventory(ctx context.Context, filter model.InventoryFilter) ([]model.Inventory, error) {
	mongoFilter := createFilter(filter)

	cursor, err := r.collection.Find(ctx, mongoFilter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, model.ErrInventoryNotFound
		}
		return nil, err
	}
	defer func() {
		cerr := cursor.Close(ctx)
		if cerr != nil {
			log.Printf("failed to close cursor: %v\n", err)
		}
	}()

	var repoInventories []repomodel.Inventory
	if err = cursor.All(ctx, &repoInventories); err != nil {
		return nil, err
	}

	Inventory := make([]model.Inventory, len(repoInventories))
	for i, repoInventory := range repoInventories {
		Inventory[i] = converter.ConvertInventoryRepoModelToModel(repoInventory)
	}
	return Inventory, nil
}

func createFilter(listpart model.InventoryFilter) bson.M {
	mongoFilter := bson.M{}
	if listpart.UUID != nil && len(*listpart.UUID) > 0 {
		mongoFilter["uuid"] = bson.M{"$in": *listpart.UUID}
	}

	if listpart.Names != nil && len(*listpart.Names) > 0 {
		mongoFilter["inventoryinfo.name"] = bson.M{"$in": *listpart.Names}
	}

	if listpart.Categories != nil && len(*listpart.Categories) > 0 {
		mongoFilter["inventoryinfo.category"] = bson.M{"$in": *listpart.Categories}
	}

	if listpart.ManufacturerCountries != nil && len(*listpart.ManufacturerCountries) > 0 {
		mongoFilter["nventoryinfo.manufacturer.country"] = bson.M{"$in": *listpart.ManufacturerCountries}
	}

	if listpart.Tags != nil && len(*listpart.Tags) > 0 {
		mongoFilter["inventoryinfo.tags"] = bson.M{"$in": *listpart.Tags}
	}
	return mongoFilter
}
