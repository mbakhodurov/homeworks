package payment

import (
	"context"

	"github.com/google/uuid"
	"github.com/mbakhodurov/homeworks/week7/payment/internal/service"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *Service) PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod service.PaymentMethod) (string, error) {
	transactionUUID := uuid.NewString()

	logger.Info(ctx, "payment succeeded",
		zap.String("transaction_uuid", transactionUUID),
		zap.String("order_uuid", orderUUID),
		zap.String("user_uuid", userUUID),
		zap.Int("payment_method", int(paymentMethod)),
	)

	return transactionUUID, nil
}
