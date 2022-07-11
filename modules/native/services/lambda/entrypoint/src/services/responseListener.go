package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/olebedev/emitter"
	amqp "github.com/rabbitmq/amqp091-go"
)

func ListenForResponces(ctx context.Context, amqpClient *amqp.Channel, emiter *emitter.Emitter) (string, context.CancelFunc, error) {
	queue, err := amqpClient.QueueDeclare("", false, true, true, false, amqp.Table{})
	if err != nil {
		return "", nil, err
	}

	consumeCh, err := amqpClient.Consume(queue.Name, "response_receiver", true, false, false, false, amqp.Table{})
	if err != nil {
		return "", nil, err
	}

	replyKey := uuid.New().String()
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		for {
			select {
			case delivery := <-consumeCh:
				done := emiter.Emit(delivery.CorrelationId, delivery)
				select {
				case <-done:
				case <-time.After(time.Millisecond * 50):
					close(done)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return replyKey, cancel, nil
}
