package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks/week3/order/internal/client/converter"
	"github.com/mbakhodurov/homeworks/week3/order/internal/model"
	payment_v1 "github.com/mbakhodurov/homeworks/week3/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, paymentMethod model.PaymentMethod, orderUUID, userUUID string) (string, error) {
	res, err := c.generatedClient.PayOrder(ctx, &payment_v1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: converter.PaymentMethodModelToPaymentProtoPaymentModel(paymentMethod),
	})

	if err != nil {
		return "", err
	}

	return res.TransactionUuid, nil
}
