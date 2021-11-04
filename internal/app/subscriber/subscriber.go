package subscriber

import (
	"context"
	"errors"

	"github.com/italolelis/fourkeys/internal/app/subscriber/awslambda"
	"github.com/italolelis/fourkeys/internal/app/subscriber/kinesis"
)

type (
	// Subscriber is the interface that must be implemented by the subscriber.
	Subscriber interface {
		Subscribe(ctx context.Context, streamName string, fn func(ctx context.Context, payload []byte, headers map[string][]string) error) error
	}
)

func Build(ctx context.Context, driver string, configs map[string]string) (Subscriber, error) {
	switch driver {
	case "awslambda":
		return awslambda.NewSubscriber(ctx)
	case "kinesis":
		return kinesis.NewSubscriber(ctx, kinesis.SessionConfig{
			Endpoint:   configs["endpoint"],
			Region:     configs["region"],
			StreamName: configs["stream_name"],
		})
	default:
		return nil, errors.New("driver not supported. Please use one of the supported ones: awslambda, or kinesis")
	}
}
