package inventory

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *repository) Update(ctx context.Context, uuid string, set bson.M) error {
	if len(set) == 0 {
		return nil
	}

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"uuid": uuid},
		bson.M{"$set": set},
	)
	if err != nil {
		return err
	}

	return nil
}
