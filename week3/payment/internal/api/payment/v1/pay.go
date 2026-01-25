package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks/week3/payment/internal/service"
	payment_v1 "github.com/mbakhodurov/homeworks/week3/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) PayOrder(ctx context.Context, req *payment_v1.PayOrderRequest) (*payment_v1.PayOrderResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	orderUUID := req.GetOrderUuid()
	userUUID := req.GetUserUuid()
	paymentMethod := req.GetPaymentMethod()

	transactionUUID, err := a.service.PayOrder(ctx, orderUUID, userUUID, service.PaymentMethod(paymentMethod))
	if err != nil {
		return nil, err
	}

	return &payment_v1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
