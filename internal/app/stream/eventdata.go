package stream

import (
	"context"
	"fmt"

	"github.com/italolelis/watchops/internal/app/provider"
)

type (
	// EventDataHandler is a handler for event data.
	EventDataHandler struct {
		w  provider.Writer
		pr *provider.ParserRegistry
	}

	// Handler is the interface that must be implemented by the handler.
	HandlerFunc func(ctx context.Context, payload []byte, headers map[string][]string) error
)

// NewEventDataHandler create a new handler for event data.
func NewEventDataHandler(w provider.Writer, pr *provider.ParserRegistry) *EventDataHandler {
	return &EventDataHandler{w: w, pr: pr}
}

// Handle handles all events that come.
func (h *EventDataHandler) Handle(ctx context.Context, payload []byte, headers map[string][]string) error {
	source := headers["source"]
	if len(source) == 0 {
		return fmt.Errorf("source header is required")
	}

	p := h.pr.Get(source[0])
	eventData, err := p.Parse(headers, payload)
	if err != nil {
		return fmt.Errorf("failed to parse payload from webhook: %w", err)
	}

	if id, ok := headers["msg_id"]; ok {
		eventData.MsgID = id[0]
	}

	if err := h.w.Add(ctx, eventData); err != nil {
		return fmt.Errorf("failed to store event in data store: %w", err)
	}

	return nil
}
