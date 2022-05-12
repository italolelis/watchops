package pubsub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"cloud.google.com/go/pubsub"
	"github.com/italolelis/watchops/internal/app/wh"
	"github.com/italolelis/watchops/internal/pkg/log"
	"github.com/italolelis/watchops/internal/pkg/signal"
)

var (
	ErrSubscriptionDoesNotExists = errors.New("subscription does not exists")
)

type SessionConfig struct {
	ProjectID    string `split_words:"true"`
	Subscription string `split_words:"true"`
}

// Subscriber is the kinesis subscriber.
type Subscriber struct {
	client         *pubsub.Client
	SubscriptionID string
}

// NewSubscriber creates a new instance of Subscriber.
func NewSubscriber(ctx context.Context, cfg SessionConfig) (*Subscriber, error, func() error) {
	logger := log.WithContext(ctx).Named("pubsub_subscriber").With("project", cfg.ProjectID)

	logger.Debugw("building subscriber")
	client, err := pubsub.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err), nil
	}

	return &Subscriber{
		client:         client,
		SubscriptionID: cfg.Subscription,
	}, nil, client.Close
}

// Subscribe subscribes to the kinesis stream.
//nolint: exhaustivestruct
func (s *Subscriber) Subscribe(ctx context.Context, fn func(ctx context.Context, payload []byte, headers map[string][]string) error) error {
	logger := log.WithContext(ctx).Named("pubsub_subscriber")

	ctx, cancel := context.WithCancel(ctx)

	done := signal.New(ctx)
	go func() {
		<-done.Done()
		cancel()
	}()

	logger.Infow("getting subscription", "sub", s.SubscriptionID)
	sub := s.client.Subscription(s.SubscriptionID)

	exists, err := sub.Exists(ctx)
	if err != nil {
		return fmt.Errorf("pubsub subscription %s failed: %w", s.SubscriptionID, err)
	}

	if !exists {
		return ErrSubscriptionDoesNotExists
	}

	// Receiving messages.
	logger.Infow("processing messages...", "sub", sub)
	return sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		logger.Debug("incoming message received")

		var messageContainer wh.MessageContainer
		if err := json.Unmarshal(msg.Data, &messageContainer); err != nil {
			logger.Errorw("failed to parse message from stream", "err", err)
			return
		}

		messageContainer.Headers["msg_id"] = append(messageContainer.Headers["msg_id"], msg.ID)

		arrivalTime := msg.PublishTime
		messageContainer.Headers["publish_time"] = append(messageContainer.Headers["publish_time"], strconv.FormatInt(arrivalTime.Unix(), 10))
		messageContainer.Headers["source"] = append(messageContainer.Headers["source"], messageContainer.Source)

		if err := fn(ctx, messageContainer.Payload, messageContainer.Headers); err != nil {
			logger.Errorw("failed to process event", "err", err)
			return
		}

		msg.Ack()
		logger.Debug("finished processing incoming message")
	})
}
