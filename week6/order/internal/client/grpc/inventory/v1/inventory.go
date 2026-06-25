package v1

import (
	"context"
	"fmt"

	"github.com/mbakhodurov/homeworks/week6/order/internal/client/converter"
	"github.com/mbakhodurov/homeworks/week6/order/internal/model"
	grpcAuth "github.com/mbakhodurov/homeworks/week6/platform/pkg/middleware/grpc"
	inventory_v1 "github.com/mbakhodurov/homeworks/week6/shared/pkg/proto/inventory/v1"
)

func (c *Client) ListPartInventory(ctx context.Context, filter model.InventoryFilter) ([]model.Inventory, error) {
	// Добавляем session UUID в gRPC metadata для передачи в Inventory сервис
	ctx = grpcAuth.ForwardSessionUUIDToGRPC(ctx)
	// Для отладки:
	if uuid, ok := grpcAuth.GetSessionUUIDFromContext(ctx); ok {
		fmt.Println("Session UUID:", uuid)
	}
	if user, ok := grpcAuth.GetUserFromContext(ctx); ok {
		fmt.Println("User:", user.Info)
	}

	res, err := c.generatedClient.ListPartInventory(ctx, &inventory_v1.ListPartInventoryRequest{
		Filter: converter.InventoryFilterModelToProto(filter),
	})

	if err != nil {
		return []model.Inventory{}, err
	}

	return converter.ListPartInventoryResponseProtoToModel(res), nil
}
