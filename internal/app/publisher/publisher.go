package publisher

import "context"

// Publisher represents the message publisher interface.
type Publisher interface {
	Publish(ctx context.Context, streamName string, data []byte) error
}
