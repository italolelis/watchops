package kinesis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	awskinesis "github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
)

var (
	ErrEmptyStreamName = errors.New("kinesis stream name can't be empty")
	ErrEmptySession    = errors.New("kinesis session can't be nil")
)

// Publisher holds the kinesis connection and the stream name.
type Publisher struct {
	kinesis kinesisiface.KinesisAPI
}

// NewPublisher creates a new kinesis connection.
func NewPublisher(session client.ConfigProvider) (*Publisher, error) {
	if session == nil {
		return nil, ErrEmptySession
	}

	return &Publisher{
		kinesis: awskinesis.New(session),
	}, nil
}

// Publish publish the data to the correct stream.
//nolint: exhaustivestruct
func (p *Publisher) Publish(ctx context.Context, streamName string, data []byte) error {
	_, err := p.kinesis.PutRecordWithContext(ctx, &awskinesis.PutRecordInput{
		Data:         data,
		StreamName:   &streamName,
		PartitionKey: aws.String(time.Now().Format(time.RFC3339Nano)),
	})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
