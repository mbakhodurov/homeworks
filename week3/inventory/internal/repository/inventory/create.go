package inventory

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mbakhodurov/homeworks/week3/inventory/internal/model"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/mbakhodurov/homeworks/week3/inventory/internal/repository/converter"
	repoModel "github.com/mbakhodurov/homeworks/week3/inventory/internal/repository/model"
)

func (r *repository) Create(ctx context.Context, inventoryInfo model.InventoryInfo) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	resModel := repoModel.Inventory{
		UUID:          uuid.NewString(),
		InventoryInfo: converter.ConvertInventoryInfoModelToRepoModel(inventoryInfo),
		CreatedAt:     time.Now(),
	}

	_, err := r.collection.InsertOne(ctx, resModel)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", fmt.Errorf("inventory already exists: %v", err)
		}
		return "", err
	}

	return resModel.UUID, nil
}
