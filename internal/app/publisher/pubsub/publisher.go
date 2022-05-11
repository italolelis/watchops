package pubsub

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
)

var (
	ErrTopicDoesNotExist = errors.New("pubsub topic doesn't exist")
	ErrEmptySession      = errors.New("pubsub session can't be nil")
)

type SessionConfig struct {
	ProjectID string        `split_words:"true"`
	Timeout   time.Duration `default:"5s"`
}

// Publisher holds the kinesis connection and the stream name.
type Publisher struct {
	client  *pubsub.Client
	timeout time.Duration
}

// NewPublisher creates a new kinesis connection.
func NewPublisher(ctx context.Context, cfg SessionConfig) (*Publisher, error, func() error) {
	client, err := pubsub.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err), nil
	}

	return &Publisher{
		client:  client,
		timeout: cfg.Timeout,
	}, nil, client.Close
}

// Publish publish the data to the correct stream.
//nolint: exhaustivestruct
func (p *Publisher) Publish(ctx context.Context, topicID string, data []byte) error {
	topic := p.client.Topic(topicID)

	exists, err := topic.Exists(ctx)
	if err != nil {
		return fmt.Errorf("pubsub topic %s failed: %w", topicID, err)
	}

	if !exists {
		return ErrTopicDoesNotExist
	}

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	result := topic.Publish(ctx, &pubsub.Message{
		Data: data,
	})

	if _, err := result.Get(ctx); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
