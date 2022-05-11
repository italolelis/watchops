package subscriber

import (
	"context"
	"errors"

	"github.com/italolelis/watchops/internal/app/subscriber/awslambda"
	"github.com/italolelis/watchops/internal/app/subscriber/kinesis"
	"github.com/italolelis/watchops/internal/app/subscriber/pubsub"
)

type (
	// Subscriber is the interface that must be implemented by the subscriber.
	Subscriber interface {
		Subscribe(ctx context.Context, fn func(ctx context.Context, payload []byte, headers map[string][]string) error) error
	}

	// Config is the configuration for the subscriber.
	Config struct {
		Driver  string
		Kinesis kinesis.SessionConfig
		Pubsub  pubsub.SessionConfig
	}
)

func Build(ctx context.Context, driver string, cfg Config) (Subscriber, error, func() error) {
	switch driver {
	case "awslambda":
		s, err := awslambda.NewSubscriber(ctx)
		return s, err, func() error { return nil }
	case "kinesis":
		s, err := kinesis.NewSubscriber(ctx, cfg.Kinesis)
		return s, err, func() error { return nil }
	case "pubsub":
		return pubsub.NewSubscriber(ctx, cfg.Pubsub)
	default:
		return nil, errors.New("driver not supported. Please use one of the supported ones: awslambda, kinesis, and pubsub"), func() error { return nil }
	}
}
