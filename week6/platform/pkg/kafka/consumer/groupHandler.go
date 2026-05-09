package consumer

import (
	"github.com/IBM/sarama"
	"github.com/mbakhodurov/homeworks/week6/platform/pkg/kafka"
	"go.uber.org/zap"
)

type Middleware func(next kafka.MessageHandler) kafka.MessageHandler

type groupHandler struct {
	handler kafka.MessageHandler
	logger  Logger
}

// NewGroupHandler создаёт новый groupHandler с middleware цепочкой.
func NewGroupHandler(handler kafka.MessageHandler, logger Logger, middlewares ...Middleware) *groupHandler {
	// Применяем middleware цепочку
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return &groupHandler{
		handler: handler,
		logger:  logger,
	}
}

func (g *groupHandler) Setup(session sarama.ConsumerGroupSession) error {
	g.logger.Info(session.Context(), "New session",
		zap.String("member_id", session.MemberID()),
		zap.Any("claims", session.Claims()),
	)
	return nil
	// return nil
}

func (g *groupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (g *groupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			g.logger.Info(session.Context(), "message received",
				zap.String("member_id", session.MemberID()),
				zap.ByteString("key", message.Key),
				zap.Int32("partition", message.Partition),
				zap.Int64("offset", message.Offset),
			)
			if !ok {
				g.logger.Info(session.Context(), "Kafka message channel closed")
				return nil
			}
			msg := kafka.Message{
				Key:            message.Key,
				Value:          message.Value,
				Topic:          message.Topic,
				Partition:      message.Partition,
				Offset:         message.Offset,
				Timestamp:      message.Timestamp,
				BlockTimestamp: message.BlockTimestamp,
				Headers:        extractHeaders(message.Headers),
			}
			if err := g.handler(session.Context(), msg); err != nil {
				g.logger.Error(session.Context(), "Kafka handler error", zap.Error(err))
				continue
			}
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			g.logger.Info(session.Context(), "Kafka session context done")
			return nil
		}

	}
}

func extractHeaders(headers []*sarama.RecordHeader) map[string][]byte {
	result := make(map[string][]byte)
	for _, h := range headers {
		if h != nil && h.Key != nil {
			result[string(h.Key)] = h.Value
		}
	}

	return result
}
