package kafka

import (
	"github.com/mbakhodurov/homeworks/week6/order/internal/model"
)

type OrderAssembledDecoder interface {
	Decode(data []byte) (model.ShipAssembled, error)
}
