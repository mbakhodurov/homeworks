package order_consumer

import (
	"context"
	"math/rand"
	"time"

	"github.com/mbakhodurov/homeworks/week7/assembly/internal/model"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/kafka"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) OrderHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode order paid event", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("event_uuid", event.EventUUID),
		zap.String("payment_method", event.PaymentMethod),
		zap.String("transaction_uuid", event.TransactionUUID),
		// zap.ByteString("key", (msg.Key)),
	)

	// time.Sleep(5 * time.Second)
	//nolint:gosec
	delay := time.Duration(rand.Intn(10)+1) * time.Second
	// select {
	// case <-time.After(delay):
	// case <-ctx.Done():
	// 	return ctx.Err()
	// }

	shipAssembled := model.ShipAssembled{
		EventUUID:    event.EventUUID,
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: int64(delay / time.Second),
	}
	err = s.orderProducer.ProduceShipAssembled(ctx, shipAssembled)
	if err != nil {
		logger.Error(ctx, "Failed to produce ship assembled event",
			zap.Any("ship_assembled", shipAssembled),
			zap.Error(err),
		)
		return err
	}

	s.assemblyHistogram.Record(ctx, float64(shipAssembled.BuildTimeSec))

	return nil
}
