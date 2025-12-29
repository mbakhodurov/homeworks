package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks/week2/inventory/internal/converter"
	inventory_v1 "github.com/mbakhodurov/homeworks/week2/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (a *InventoryApi) Update(ctx context.Context, req *inventory_v1.UpdateRequest) (*emptypb.Empty, error) {
	if req.GetUpdateInfo() == nil {
		return nil, status.Error(codes.InvalidArgument, "update_info cannot be nil")
	}
	err := a.inventoryService.Update(ctx, converter.InventoryUpdateInfoToModel(req.GetUpdateInfo()), req.GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update sighting: %v", err)

	}
	return &emptypb.Empty{}, nil
}
