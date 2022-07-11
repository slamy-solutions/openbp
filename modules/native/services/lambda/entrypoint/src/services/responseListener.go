package services

import (
	"context"
	"time"

	"github.com/olebedev/emitter"
	amqp "github.com/rabbitmq/amqp091-go"
)

const RESPONSE_EXCHANGE = "native_lambda_entrypoint_output"

func ListenForResponces(ctx context.Context, amqpClient *amqp.Channel, emiter *emitter.Emitter) (string, context.CancelFunc, error) {
	err := amqpClient.ExchangeDeclare(RESPONSE_EXCHANGE, "direct", true, false, false, false, nil)
	if err != nil {
		return "", nil, err
	}

	queue, err := amqpClient.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		return "", nil, err
	}

	err = amqpClient.QueueBind(queue.Name, queue.Name, RESPONSE_EXCHANGE, false, nil)
	if err != nil {
		return "", nil, err
	}

	consumeCh, err := amqpClient.Consume(queue.Name, "response_receiver", true, false, false, false, nil)
	if err != nil {
		return "", nil, err
	}

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
	return queue.Name, cancel, nil
}
