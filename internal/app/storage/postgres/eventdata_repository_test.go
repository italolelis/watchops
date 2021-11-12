package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/italolelis/watchops/internal/app/provider"
	"github.com/stretchr/testify/require"
)

func TestEventDataRepository_Add(t *testing.T) {
	ctx := context.Background()

	db, mock := newDBMock()
	defer db.Close()

	writer := NewEventDataWriter(db)

	e := provider.Event{
		ID:          "event-id",
		EventType:   "event-name",
		Metadata:    json.RawMessage(`{"key": "value"}`),
		TimeCreated: time.Now(),
		Signature:   "signature",
		MsgID:       "msg-id",
		Source:      "source",
	}

	const query = "INSERT INTO watchops.events_raw"

	mock.ExpectExec(query).
		WithArgs(e.ID, e.EventType, e.Metadata, e.TimeCreated, e.Signature, e.MsgID, e.Source).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := writer.Add(ctx, e)
	require.NoError(t, err)
}

func TestEventDataRepository_Add_Failure(t *testing.T) {
	ctx := context.Background()

	db, mock := newDBMock()
	defer db.Close()

	writer := NewEventDataWriter(db)

	e := provider.Event{
		ID:          "event-id",
		EventType:   "event-name",
		Metadata:    json.RawMessage(`{"key": "value"}`),
		TimeCreated: time.Now(),
		Signature:   "signature",
		MsgID:       "msg-id",
		Source:      "source",
	}

	const query = "INSERT INTO watchops.events_raw"

	mock.ExpectExec(query).
		WithArgs(e.ID, e.EventType, e.Metadata, e.TimeCreated, e.Signature, e.MsgID, e.Source).
		WillReturnError(fmt.Errorf("failed to insert"))

	err := writer.Add(ctx, e)
	require.Error(t, err)
}

func TestEventDataRepository_Close(t *testing.T) {
	db, mock := newDBMock()
	mock.ExpectClose()

	require.NoError(t, db.Close())
}

func newDBMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}
