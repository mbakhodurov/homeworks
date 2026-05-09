package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks/week6/order/internal/model"
	grpcAuth "github.com/mbakhodurov/homeworks/week6/platform/pkg/middleware/grpc"
	payment_v1 "github.com/mbakhodurov/homeworks/week6/shared/pkg/proto/payment/v1"
)

func (c *Client) PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod model.PaymentMethod) (string, error) {
	ctx = grpcAuth.ForwardSessionUUIDToGRPC(ctx)

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
