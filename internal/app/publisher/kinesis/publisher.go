package kinesis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

var (
	ErrEmptyStreamName = errors.New("kinesis stream name can't be empty")
	ErrEmptySession    = errors.New("kinesis session can't be nil")
)

type SessionConfig struct {
	Endpoint   string
	Region     string
	Timeout    time.Duration `default:"5s"`
	MaxRetries int           `split_words:"true" default:"3"`
}

// Publisher holds the kinesis connection and the stream name.
type Publisher struct {
	kinesis *kinesis.Client
	timeout time.Duration
}

// NewPublisher creates a new kinesis connection.
func NewPublisher(ctx context.Context, cfg SessionConfig) (*Publisher, error) {
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

	return &Publisher{
		kinesis: kinesis.NewFromConfig(awsCfg),
		timeout: cfg.Timeout,
	}, nil
}

// Publish publish the data to the correct stream.
//nolint: exhaustivestruct
func (p *Publisher) Publish(ctx context.Context, streamName string, data []byte) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	_, err := p.kinesis.PutRecord(ctx, &kinesis.PutRecordInput{
		Data:         data,
		StreamName:   &streamName,
		PartitionKey: aws.String(time.Now().Format(time.RFC3339Nano)),
	})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
