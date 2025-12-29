package inventory

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mbakhodurov/homeworks/week2/inventory/internal/model"
	"github.com/mbakhodurov/homeworks/week2/inventory/internal/repository/converter"
	repoModel "github.com/mbakhodurov/homeworks/week2/inventory/internal/repository/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *repository) Create(ctx context.Context, part model.PartInfo) (string, error) {
	newUUID := uuid.NewString()

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, v := range r.data {
		if v.Partinfo.Name == part.Name {
			return "", status.Error(codes.AlreadyExists, "part already exists")
		}
	}

	r.data[newUUID] = repoModel.Part{
		UUID:      newUUID,
		Partinfo:  converter.InventoryInfoToRepoModel(part),
		CreatedAt: time.Now(),
	}
	return r.data[newUUID].UUID, nil
}
