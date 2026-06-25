package kafka

import (
	"github.com/mbakhodurov/homeworks/week7/assembly/internal/model"
)

type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaid, error)
}
