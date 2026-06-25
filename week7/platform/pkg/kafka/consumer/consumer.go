package consumer

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/kafka"
	"go.uber.org/zap"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type consumer struct {
	group       sarama.ConsumerGroup
	topic       []string
	logger      Logger
	middlewares []Middleware
}

func NewConsumer(group sarama.ConsumerGroup, topics []string, logger Logger, middlewares ...Middleware) *consumer {
	return &consumer{
		group:       group,
		topic:       topics,
		logger:      logger,
		middlewares: middlewares,
	}
}

func (c *consumer) Consume(ctx context.Context, handler kafka.MessageHandler) error {
	newGroupHandler := NewGroupHandler(handler, c.logger, c.middlewares...)

	for {
		if err := c.group.Consume(ctx, c.topic, newGroupHandler); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return nil
			}

			c.logger.Error(ctx, "Kafka consume error", zap.Error(err))
			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}
		c.logger.Info(ctx, "Kafka consumer group rebalancing...")
	}
}
