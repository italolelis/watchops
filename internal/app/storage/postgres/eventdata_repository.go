package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/italolelis/fourkeys/internal/app/provider"
)

type (
	// EventDataWriter implement writer interface.
	EventDataWriter struct{ db *sql.DB }
)

// NewEventDataWriter creates a new instance of EventDataWriter.
func NewEventDataWriter(db *sql.DB) *EventDataWriter {
	return &EventDataWriter{db: db}
}

// Add adds event data coming from webhooks.
func (w *EventDataWriter) Add(ctx context.Context, eventData provider.Event) error {
	const query = `
		INSERT INTO 
			fourkeys.events_raw(id, event_type, metadata, time_created, signature, msg_id, source) 
			VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	if _, err := w.db.ExecContext(
		ctx,
		query,
		eventData.ID,
		eventData.EventType,
		eventData.Metadata,
		eventData.TimeCreated,
		eventData.Signature,
		eventData.MsgID,
		eventData.Source,
	); err != nil {
		return fmt.Errorf("error when executing event data query: %w", err)
	}

	return nil
}

// Closes the connection to the database.
func (w *EventDataWriter) Close() error {
	return w.db.Close()
}
