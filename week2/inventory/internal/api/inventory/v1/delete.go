package v1

import (
	"context"
	"errors"

	"github.com/mbakhodurov/homeworks/week2/inventory/internal/model"
	inventory_v1 "github.com/mbakhodurov/homeworks/week2/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (a *InventoryApi) DeletePart(ctx context.Context, req *inventory_v1.DeletePartRequest) (*emptypb.Empty, error) {
	err := a.inventoryService.DeleteByUUID(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, model.ErrInventoryNotFound) {
			return nil, status.Errorf(codes.NotFound, "sighting with UUID %s not found", req.GetUuid())
		}
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
