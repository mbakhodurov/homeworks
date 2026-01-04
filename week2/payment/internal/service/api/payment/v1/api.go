package v1

import (
	"github.com/mbakhodurov/homeworks/week2/payment/internal/service"
	payment_v1 "github.com/mbakhodurov/homeworks/week2/shared/pkg/proto/payment/v1"
)

type api struct {
	payment_v1.UnimplementedPaymentServiceServer

	paymentService service.PaymentService
}

func NewApi(paymentService service.PaymentService) *api {
	return &api{
		paymentService: paymentService,
	}
}
