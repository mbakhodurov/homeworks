package inventory

import (
	"context"
	"sync"

	def "github.com/mbakhodurov/homeworks/week3/inventory/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	mu         sync.RWMutex
	collection *mongo.Collection
}

func NewRepository(ctx context.Context, db *mongo.Database) *repository {
	collection := db.Collection("inventories")

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "inventory_info.name", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "uuid", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		panic(err)
	}

	return &repository{
		collection: collection,
	}
}
