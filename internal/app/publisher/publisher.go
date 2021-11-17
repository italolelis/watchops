package publisher

import (
	"context"
	"errors"

	"github.com/italolelis/watchops/internal/app/publisher/kinesis"
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
	}
)

func Build(ctx context.Context, cfg Config) (Publisher, error) {
	switch cfg.Driver {
	case "kinesis":
		return kinesis.NewPublisher(ctx, cfg.Kinesis)
	default:
		return nil, errors.New("driver not supported. Please use one of the supported ones: kinesis")
	}
}
