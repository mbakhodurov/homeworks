package order_producer

import (
	def "github.com/mbakhodurov/homeworks/week5/order/internal/service"
	"github.com/mbakhodurov/homeworks/week5/platform/pkg/kafka"
)

var _ def.OrderProducerService = (*service)(nil)

type service struct {
	orderProducer kafka.Producer
}

func NewService(orderProducer kafka.Producer) *service {
	return &service{
		orderProducer: orderProducer,
	}
}
