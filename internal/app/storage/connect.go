package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/italolelis/watchops/internal/app/provider"
	"github.com/italolelis/watchops/internal/app/storage/bigquery"
	"github.com/italolelis/watchops/internal/app/storage/postgres"
	"github.com/italolelis/watchops/internal/app/storage/redshift"
	_ "github.com/lib/pq"
)

// ErrInvalidDataSource is returned when the data source is not supported.
var ErrInvalidDataSource = errors.New("invalid data source")

type Config struct {
	Driver          string
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	Timeout         time.Duration
	SchemaName      string

	Bigquery struct {
		ProjectID string
	}
}

// Connect creates a new connection to the database. It returns a provider Writer based
// on the given configuration.
func Connect(ctx context.Context, cfg Config) (provider.Writer, error) {
	switch cfg.Driver {
	case "postgres":
		db, err := sql.Open("postgres", cfg.DSN)
		if err != nil {
			return nil, fmt.Errorf("could not connect to the database: %w", err)
		}

		db.SetMaxOpenConns(cfg.MaxOpenConns)
		db.SetMaxIdleConns(cfg.MaxIdleConns)
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

		return postgres.NewEventDataWriter(db, cfg.SchemaName), nil
	case "redshift":
		db, err := sql.Open("postgres", cfg.DSN)
		if err != nil {
			return nil, fmt.Errorf("could not connect to the database: %w", err)
		}

		db.SetMaxOpenConns(cfg.MaxOpenConns)
		db.SetMaxIdleConns(cfg.MaxIdleConns)
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

		return redshift.NewEventDataWriter(db, cfg.SchemaName), nil
	case "bigquery":
		return bigquery.NewEventDataWriter(ctx, cfg.Bigquery.ProjectID)
	default:
		return nil, ErrInvalidDataSource
	}
}
