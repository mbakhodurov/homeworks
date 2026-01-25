package payment

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/mbakhodurov/homeworks/week3/payment/internal/service"
)

func (s *Service) PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod service.PaymentMethod) (string, error) {
	transactionUUID := uuid.NewString()
	log.Printf("Оплата прошла успешно, transactionUUID: %s, orderUUID: %s, userUUID: %s, paymentMethod: %d",
		transactionUUID, orderUUID, userUUID, paymentMethod)
	return transactionUUID, nil
}
