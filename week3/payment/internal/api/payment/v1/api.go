package v1

import (
	"github.com/mbakhodurov/homeworks/week3/payment/internal/service"
	payment_v1 "github.com/mbakhodurov/homeworks/week3/shared/pkg/proto/payment/v1"
)

type api struct {
	service service.PaymentService
	payment_v1.UnimplementedPaymentServiceServer
}

func NewApi(service service.PaymentService) *api {
	return &api{
		service: service,
	}
}
