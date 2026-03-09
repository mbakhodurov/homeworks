package payment

import (
	def "github.com/mbakhodurov/homeworks/week4/payment/internal/service"
)

var _ def.PaymentService = (*Service)(nil)

type Service struct{}

func NewService() *Service {
	return &Service{}
}
