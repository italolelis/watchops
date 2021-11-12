package awslambda

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/italolelis/watchops/internal/app/wh"
	"github.com/italolelis/watchops/internal/pkg/log"
)

type handlerFunc func(ctx context.Context, payload []byte, headers map[string][]string) error

// Subscriber is the lambda kinesis subscriber.
type Subscriber struct{}

// NewSubscriber creates a new instance of Subscriber.
func NewSubscriber(ctx context.Context) (*Subscriber, error) {
	return &Subscriber{}, nil
}

// Subscribe subscribes to the kinesis stream for a lambda. This method will block the caller.
//nolint: exhaustivestruct
func (s *Subscriber) Subscribe(ctx context.Context, fn func(ctx context.Context, payload []byte, headers map[string][]string) error) error {
	lambda.StartWithContext(ctx, s.lambdaHandler(fn))

	return nil
}

func (s *Subscriber) lambdaHandler(fn handlerFunc) func(context.Context, events.KinesisEvent) error {
	return func(ctx context.Context, kinesisEvent events.KinesisEvent) error {
		logger := log.WithContext(ctx).Named("lambda_subscriber")
		logger.Info("processing messages...")

		for _, record := range kinesisEvent.Records {
			logger.Debug("incoming message received")

			data := record.Kinesis.Data

			var messageContainer wh.MessageContainer
			if err := json.Unmarshal(data, &messageContainer); err != nil {
				return fmt.Errorf("failed to parse message from stream: %w", err)
			}

			messageContainer.Headers["msg_id"] = append(messageContainer.Headers["msg_id"], record.EventID)

			if err := fn(ctx, messageContainer.Payload, messageContainer.Headers); err != nil {
				logger.Errorw("failed to process event", "err", err)

				continue
			}

			logger.Debug("finished processing incoming message")
		}

		logger.Info("finished processing messages")

		return nil
	}
}
