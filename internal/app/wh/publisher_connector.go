// Package wh is the webhook manager writer. The goal of this package is
// to manage an input (webhook of any kind), send it to a parser, and write
// to an output.
package wh

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/italolelis/watchops/internal/app/publisher"
)

// ErrSourceNotSupported is returned when the source is not supported.
var ErrSourceNotSupported = errors.New("source not supported")

// Connector holds all connector methods.
type Connector interface {
	Write(ctx context.Context, payload []byte, headers map[string][]string) error
}

// PublisherConnector is a connector for the publisher.
type PublisherConnector struct {
	pub         publisher.Publisher
	topicPrefix string
}

// MessageContainer is a container for the message that will be sent to the publisher.
type MessageContainer struct {
	Source  string              `json:"source"`
	Payload []byte              `json:"payload"`
	Headers map[string][]string `json:"headers"`
}

// NewPublisherConnector creates a new instance of PublisherConnector.
func NewPublisherConnector(ctx context.Context, p publisher.Publisher, topicPrefix string) *PublisherConnector {
	return &PublisherConnector{pub: p, topicPrefix: topicPrefix}
}

// Write sends the payload to the publisher.
func (w *PublisherConnector) Write(ctx context.Context, payload []byte, headers map[string][]string) error {
	source := getSource(headers)
	if source == "" {
		return ErrSourceNotSupported
	}

	raw, err := json.Marshal(MessageContainer{
		Source:  source,
		Payload: payload,
		Headers: headers,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal raw data to send to the message broker: %w", err)
	}

	return w.pub.Publish(ctx, w.topicPrefix+source, raw)
}

func getSource(headers map[string][]string) string {
	source := strings.TrimSpace(strings.Split(headers["User-Agent"][0], "/")[0])

	// This will eventually grow. So we leave it as a switch case for now.
	switch source {
	case "GitHub-Hookshot":
		return "github"
	case "Opsgenie Http Client":
		return "opsgenie"
	case "X-Gitlab-Event":
		return "gitlab"
	case "Circleci-Event-Type":
		return "circleci"
	default:
		return ""
	}
}
