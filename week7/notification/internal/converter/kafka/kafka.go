package kafka

import "github.com/mbakhodurov/homeworks/week7/notification/internal/model"

type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaid, error)
}

type OrderAssembledDecoder interface {
	Decode(data []byte) (model.ShipAssembled, error)
}
