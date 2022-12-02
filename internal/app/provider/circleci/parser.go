package circleci

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/italolelis/watchops/internal/app/provider"
)

const (
	sourceType      = "circleci"
	eventTypeHeader = "Circleci-Event-Type"
	signatureHeader = "Circleci-Signature"
)

type event struct {
	ID         string    `json:"id"`
	HappenedAt time.Time `json:"happened_at"`
}

// Parser is the GitHub webhook parser.
type Parser struct{}

// GetName returns the parser name.
func (p *Parser) GetName() string {
	return sourceType
}

// Parse parses the incoming webhook.
func (p *Parser) Parse(headers map[string][]string, payload []byte) (provider.Event, error) {
	var e event
	if err := json.Unmarshal(payload, &e); err != nil {
		return provider.Event{}, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	event := provider.Event{
		EventType: headers[eventTypeHeader][0],
		Signature: headers[signatureHeader][0],
		Source:    sourceType,
		Metadata:  payload,
		MsgID:     headers["msg_id"][0],
	}

	switch event.EventType {
	case "workflow-completed",
		"job-completed":
		event.ID = e.ID
		event.TimeCreated = e.HappenedAt
	default:
		return provider.Event{}, &provider.UnkownTypeError{Type: event.EventType}
	}

	return event, nil
}
