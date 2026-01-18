package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks/week3/inventory/internal/converter"
	inventory_v1 "github.com/mbakhodurov/homeworks/week3/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (a *InventoryApi) UpdateInventory(ctx context.Context, req *inventory_v1.UpdateInventoryRequest) (*emptypb.Empty, error) {
	if req.GetUpdateInfo() == nil {
		return nil, status.Error(codes.InvalidArgument, "update_info cannot be nil")
	}

	if err := a.inventoryService.Update(ctx, req.GetUuid(), converter.InventoryUpdateInfoProtoToModel(req.GetUpdateInfo())); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update inventory:%v", err)
	}

	return &emptypb.Empty{}, nil
}
