package app

import (
	"context"

	v1 "github.com/mbakhodurov/homeworks/week5/payment/internal/api/payment/v1"
	"github.com/mbakhodurov/homeworks/week5/payment/internal/service"
	paymentService "github.com/mbakhodurov/homeworks/week5/payment/internal/service/payment"
	payment_v1 "github.com/mbakhodurov/homeworks/week5/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentService service.PaymentService
	paymentV1Api   payment_v1.PaymentServiceServer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (di *diContainer) PaymentV1Api(ctx context.Context) payment_v1.PaymentServiceServer {
	if di.paymentV1Api == nil {
		di.paymentV1Api = v1.NewApi(di.PaymentService(ctx))
	}
	return di.paymentV1Api
}

func (di *diContainer) PaymentService(ctx context.Context) service.PaymentService {
	if di.paymentService == nil {
		di.paymentService = paymentService.NewService()
	}
	return di.paymentService
}
