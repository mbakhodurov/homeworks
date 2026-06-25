package decoder

import (
	"fmt"

	"github.com/mbakhodurov/homeworks/week7/order/internal/model"
	events_v1 "github.com/mbakhodurov/homeworks/week7/shared/pkg/proto/events/v1"
	"google.golang.org/protobuf/proto"
)

type decoder struct{}

func NewOrderAssembledDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.ShipAssembled, error) {
	var pb events_v1.ShipAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.ShipAssembled{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.ShipAssembled{
		EventUUID:    pb.EventUuid,
		OrderUUID:    pb.OrderUuid,
		UserUUID:     pb.UserUuid,
		BuildTimeSec: pb.BuildTimeSec,
	}, nil
}
