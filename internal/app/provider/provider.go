package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// Parser holds the methods for a webhook parser.
type Parser interface {
	GetName() string
	Parse(headers map[string][]string, payload []byte) (Event, error)
}

type SignatureValidator interface {
	Validate(r *http.Request) error
	IsSupported(incomingProvider string) bool
}

// Writer represents the writer repository.
type Writer interface {
	Add(context.Context, Event) error
	Close() error
}

// Event represents a generic data event.
type Event struct {
	EventType   string          `json:"event_type"`
	ID          string          `json:"id"`
	Metadata    json.RawMessage `json:"metadata"`
	TimeCreated time.Time       `json:"time_created"`
	Signature   string          `json:"signature"`
	MsgID       string          `json:"msg_id"`
	Source      string          `json:"source"`
}
