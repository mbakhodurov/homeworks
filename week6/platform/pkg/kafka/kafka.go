package kafka

import "context"

type MessageHandler func(ctx context.Context, msg Message) error

type Producer interface {
	Send(ctx context.Context, key, value []byte) error
}

type Consumer interface {
	Consume(ctx context.Context, handler MessageHandler) error
}
