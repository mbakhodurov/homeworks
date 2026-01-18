package inventory

import (
	"context"
	"log"

	"github.com/mbakhodurov/homeworks/week3/inventory/internal/repository/converter"
	"github.com/mbakhodurov/homeworks/week3/inventory/internal/repository/model"
	repomodel "github.com/mbakhodurov/homeworks/week3/inventory/internal/repository/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *repository) ListParts(ctx context.Context, listpart model.ListParts) (*[]model.Inventory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	mongoFilter := createFilter(listpart)

	cursor, err := r.collection.Find(ctx, mongoFilter)
	if err != nil {
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

	parts := make([]model.Inventory, len(repoInventories))

	for i, inventoryPart := range  {
		parts[i] = converter.ConvertInventoryRepoModelToModel(inventoryPart)
	}
}

func createFilter(listpart model.ListParts) bson.M {
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
