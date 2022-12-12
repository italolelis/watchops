package pagerduty

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/italolelis/watchops/internal/app/provider"
)

const (
	sourceType = "pagerduty"
)

type event struct {
	Event struct {
		ID           string    `json:"id"`
		EventType    string    `json:"event_type"`
		ResourceType string    `json:"resource_type"`
		OccurredAt   time.Time `json:"occurred_at"`
	} `json:"event"`
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
		EventType: e.Event.EventType,
		Signature: provider.GenerateSignature(payload),
		Source:    sourceType,
		Metadata:  payload,
		MsgID:     headers["msg_id"][0],
	}

	switch event.EventType {
	case "incident.triggered",
		"incident.resolved":
		event.ID = e.Event.ID
		event.TimeCreated = e.Event.OccurredAt
	default:
		return provider.Event{}, &provider.UnkownTypeError{Type: event.EventType}
	}

	return event, nil
}
