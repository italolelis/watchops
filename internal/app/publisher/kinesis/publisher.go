package kinesis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

var (
	ErrEmptyStreamName = errors.New("kinesis stream name can't be empty")
	ErrEmptySession    = errors.New("kinesis session can't be nil")
)

type SessionConfig struct {
	Endpoint string
	Region   string
}

// Publisher holds the kinesis connection and the stream name.
type Publisher struct {
	kinesis *kinesis.Client
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
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	return &Publisher{
		kinesis: kinesis.NewFromConfig(awsCfg),
	}, nil
}

// Publish publish the data to the correct stream.
//nolint: exhaustivestruct
func (p *Publisher) Publish(ctx context.Context, streamName string, data []byte) error {
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
