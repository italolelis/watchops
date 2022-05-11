package kinesis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	consumer "github.com/harlow/kinesis-consumer"
	"github.com/italolelis/watchops/internal/app/wh"
	"github.com/italolelis/watchops/internal/pkg/log"
	"github.com/italolelis/watchops/internal/pkg/signal"
)

var (
	ErrEmptyStreamName = errors.New("kinesis stream name can't be empty")
)

type SessionConfig struct {
	Endpoint   string
	Region     string
	StreamName string `split_words:"true" default:"watchops"`
	Store      StoreConfig
	Timeout    time.Duration `default:"5s"`
	MaxRetries int           `split_words:"true" default:"3"`
}

// Subscriber is the kinesis subscriber.
type Subscriber struct {
	consumer *consumer.Consumer
	timeout  time.Duration
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
		config.WithRetryer(func() aws.Retryer {
			return retry.AddWithMaxAttempts(retry.NewStandard(), cfg.MaxRetries)
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	logger.Debugw("building subscriber store", "store", cfg.Store.Driver)
	store, err := BuildStore(cfg.Store)
	if err != nil {
		return nil, fmt.Errorf("failed to build a subscriber store: %w", err)
	}

	c, err := consumer.New(
		cfg.StreamName,
		consumer.WithClient(kinesis.NewFromConfig(awsCfg)),
		consumer.WithStore(store),
		consumer.WithLogger(&LoggerAdapter{logger}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscriber: %w", err)
	}

	return &Subscriber{
		consumer: c,
		timeout:  cfg.Timeout,
	}, nil
}

// Subscribe subscribes to the kinesis stream.
//nolint: exhaustivestruct
func (s *Subscriber) Subscribe(ctx context.Context, fn func(ctx context.Context, payload []byte, headers map[string][]string) error) error {
	logger := log.WithContext(ctx).Named("kinesis_subscriber")

	ctx, cancel := context.WithCancel(ctx)

	done := signal.New(ctx)
	go func() {
		<-done.Done()
		cancel()
	}()

	logger.Info("processing messages...")
	return s.consumer.Scan(ctx, func(r *consumer.Record) error {
		logger.Debug("incoming message received")

		var messageContainer wh.MessageContainer
		if err := json.Unmarshal(r.Data, &messageContainer); err != nil {
			return fmt.Errorf("failed to parse message from stream: %w", err)
		}

		messageContainer.Headers["msg_id"] = append(messageContainer.Headers["msg_id"], *r.SequenceNumber)

		arrivalTime := *r.ApproximateArrivalTimestamp
		messageContainer.Headers["publish_time"] = append(messageContainer.Headers["publish_time"], strconv.FormatInt(arrivalTime.Unix(), 10))
		messageContainer.Headers["source"] = append(messageContainer.Headers["source"], messageContainer.Source)

		if err := fn(ctx, messageContainer.Payload, messageContainer.Headers); err != nil {
			logger.Errorw("failed to process event", "err", err)

			return nil
		}

		logger.Debug("finished processing incoming message")

		return nil
	})
}
