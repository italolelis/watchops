package kinesis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	consumer "github.com/harlow/kinesis-consumer"
	"github.com/italolelis/fourkeys/internal/app/wh"
	"github.com/italolelis/fourkeys/internal/pkg/log"
)

var (
	ErrEmptyStreamName = errors.New("kinesis stream name can't be empty")
)

type SessionConfig struct {
	Endpoint   string
	Region     string `required:"true"`
	StreamName string `split_words:"true" required:"true"`
	Store      StoreConfig
}

// Subscriber is the kinesis subscriber.
type Subscriber struct {
	client *kinesis.Client
	store  consumer.Store
}

// NewSubscriber creates a new instance of Subscriber.
func NewSubscriber(ctx context.Context, cfg SessionConfig) (*Subscriber, error) {
	logger := log.WithContext(ctx).Named("kinesis_subscriber")

	resolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		if cfg.Endpoint != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           cfg.Endpoint,
				SigningRegion: cfg.Region,
			}, nil
		}

		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	awsCfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(cfg.Region),
		config.WithEndpointResolver(resolver),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	logger.Debugw("building subscriber store", "store", cfg.Store.Driver)
	store, err := BuildStore(cfg.Store)
	if err != nil {
		return nil, fmt.Errorf("failed to build a subscriber store: %w", err)
	}

	return &Subscriber{
		client: kinesis.NewFromConfig(awsCfg),
		store:  store,
	}, nil
}

// Subscribe subscribes to the kinesis stream.
//nolint: exhaustivestruct
func (s *Subscriber) Subscribe(ctx context.Context, streamName string, fn func(ctx context.Context, payload []byte, headers map[string][]string) error) error {
	logger := log.WithContext(ctx).Named("kinesis_subscriber")

	c, err := consumer.New(
		streamName,
		consumer.WithClient(s.client),
		consumer.WithStore(s.store),
	)
	if err != nil {
		return fmt.Errorf("failed to create subscriber: %w", err)
	}

	logger.Info("processing messages...")

	return c.Scan(ctx, func(r *consumer.Record) error {
		logger.Debug("incoming message received")

		var messageContainer wh.MessageContainer
		if err := json.Unmarshal(r.Data, &messageContainer); err != nil {
			return fmt.Errorf("failed to parse message from stream: %w", err)
		}

		messageContainer.Headers["msg_id"] = append(messageContainer.Headers["msg_id"], *r.SequenceNumber)

		if err := fn(ctx, messageContainer.Payload, messageContainer.Headers); err != nil {
			logger.Errorw("failed to process event", "err", err)

			return nil
		}

		logger.Debug("finished processing incoming message")

		return nil
	})
}
