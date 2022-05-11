package publisher

import (
	"context"
	"errors"

	"github.com/italolelis/watchops/internal/app/publisher/kinesis"
	"github.com/italolelis/watchops/internal/app/publisher/pubsub"
)

type (
	// Publisher represents the message publisher interface.
	Publisher interface {
		Publish(ctx context.Context, streamName string, data []byte) error
	}

	// Config is the configuration for the subscriber.
	Config struct {
		Driver  string `required:"true"`
		Kinesis kinesis.SessionConfig
		Pubsub  pubsub.SessionConfig
	}
)

func Build(ctx context.Context, cfg Config) (Publisher, error, func() error) {
	switch cfg.Driver {
	case "kinesis":
		p, err := kinesis.NewPublisher(ctx, cfg.Kinesis)
		return p, err, func() error { return nil }
	case "pubsub":
		return pubsub.NewPublisher(ctx, cfg.Pubsub)
	default:
		return nil, errors.New("driver not supported. Please use one of the supported ones: kinesis, pubsub"), nil
	}
}
