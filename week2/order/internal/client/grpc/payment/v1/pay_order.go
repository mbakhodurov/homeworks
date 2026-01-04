package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks/week2/order/internal/model"
	payment_v1 "github.com/mbakhodurov/homeworks/week2/shared/pkg/proto/payment/v1"
)

func convertPaymentMethodToProto(method model.PaymentMethod) payment_v1.PaymentMethod {
	switch method {
	case model.PaymentMethodCard:
		return payment_v1.PaymentMethod_CARD
	case model.PaymentMethodSBP:
		return payment_v1.PaymentMethod_SBP
	case model.PaymentMethodCreditCard:
		return payment_v1.PaymentMethod_CREDIT_CARD
	case model.PaymentMethodInvestorMoney:
		return payment_v1.PaymentMethod_INVESTOR_MONEY
	default:
		return payment_v1.PaymentMethod_UNKNOWN
	}
}

func (c *client) PayOrder(ctx context.Context, orderUUID string, userUUID string, paymentMethod model.PaymentMethod) (string, error) {
	payOrderResp, err := c.generatedClient.PayOrder(ctx, &payment_v1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: convertPaymentMethodToProto(paymentMethod),
	})
	if err != nil {
		return "", err
	}

	return payOrderResp.TransactionUuid, nil
}
