package bigquery

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/italolelis/watchops/internal/app/provider"
)

const (
	datasetID = "watchops"
	tableID   = "events_raw"
)

type (
	// EventDataWriter implement writer interface.
	EventDataWriter struct {
		client *bigquery.Client
	}

	bqEvent struct {
		EventType   string
		ID          string
		Metadata    json.RawMessage
		TimeCreated time.Time
		Signature   string
		MsgID       string
		Source      string
	}
)

func (e *bqEvent) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"event_type":   e.EventType,
		"id":           e.ID,
		"metadata":     e.Metadata,
		"time_created": e.TimeCreated,
		"signature":    e.Signature,
		"msg_id":       e.MsgID,
		"source":       e.Source,
	}, bigquery.NoDedupeID, nil
}

// NewEventDataWriter creates a new instance of EventDataWriter.
func NewEventDataWriter(ctx context.Context, projectID string) (*EventDataWriter, error) {
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("could not create data operations client: %w", err)
	}

	return &EventDataWriter{client: client}, nil
}

// Add adds event data coming from webhooks.
func (w *EventDataWriter) Add(ctx context.Context, eventData provider.Event) error {
	inserter := w.client.Dataset(datasetID).Table(tableID).Inserter()

	return inserter.Put(ctx, bqEvent(eventData))
}

// Closes the connection to the database.
func (w *EventDataWriter) Close() error {
	return w.client.Close()
}
