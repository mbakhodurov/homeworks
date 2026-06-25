package v1

import (
	"context"
	"fmt"

	"github.com/mbakhodurov/homeworks/week7/order/internal/model"
	grpcAuth "github.com/mbakhodurov/homeworks/week7/platform/pkg/middleware/grpc"
	payment_v1 "github.com/mbakhodurov/homeworks/week7/shared/pkg/proto/payment/v1"
)

func (c *Client) PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod model.PaymentMethod) (string, error) {
	// Добавляем session UUID в gRPC metadata для передачи в Inventory сервис
	ctx = grpcAuth.ForwardSessionUUIDToGRPC(ctx)
	// Для отладки:
	if uuid, ok := grpcAuth.GetSessionUUIDFromContext(ctx); ok {
		fmt.Println("Session UUID:", uuid)
	}
	if user, ok := grpcAuth.GetUserFromContext(ctx); ok {
		fmt.Println("User:", user.Info)
	}

	res, err := c.generatedClient.PayOrder(ctx, &payment_v1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: payment_v1.PaymentMethod(paymentMethod),
	})

	if err != nil {
		return "", err
	}

	return res.TransactionUuid, nil
}
