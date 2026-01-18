package inventory

import (
	"context"
	"errors"
	"log"

	"github.com/mbakhodurov/homeworks/week3/inventory/internal/model"
	"github.com/mbakhodurov/homeworks/week3/inventory/internal/repository/converter"
	repoModel "github.com/mbakhodurov/homeworks/week3/inventory/internal/repository/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *repository) GetAll(ctx context.Context) ([]model.Inventory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, model.ErrInventoryNotFound
		}
		return nil, err
	}

	defer func() {
		cerr := cursor.Close(ctx)
		if cerr != nil {
			log.Printf("failed to close cursor: %v\n", cerr)
		}
	}()

	var inventories []model.Inventory
	// if err := cursor.All(ctx, &inventories); err != nil {
	// 	return nil, err
	// }

	for cursor.Next(ctx) {
		var repoInv repoModel.Inventory

		if err := cursor.Decode(&repoInv); err != nil {
			return nil, err
		}
		inventories = append(inventories, converter.ConvertInventoryRepoModelToModel(repoInv))
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(inventories) == 0 {
		return nil, model.ErrInventoryNotFound
	}

	return inventories, nil
}
