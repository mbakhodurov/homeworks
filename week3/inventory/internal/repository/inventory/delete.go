package inventory

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *repository) Delete(ctx context.Context, uuid string) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	res, err := r.collection.DeleteOne(ctx, bson.M{"uuid": uuid})
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}
