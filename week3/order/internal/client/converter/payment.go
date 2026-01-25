package converter

import (
	"github.com/mbakhodurov/homeworks/week3/order/internal/model"
	payment_v1 "github.com/mbakhodurov/homeworks/week3/shared/pkg/proto/payment/v1"
)

func PaymentMethodModelToPaymentProtoPaymentModel(from model.PaymentMethod) payment_v1.PaymentMethod {
	switch from {
	case model.PaymentMethodCard:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSBP:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCreditCard:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodInvestorMoney:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_UNKNOWN_UNSPECIFIED
	}
}
