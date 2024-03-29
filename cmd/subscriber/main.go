package main

import (
	"context"
	"fmt"
	"time"

	"github.com/italolelis/watchops/internal/app/provider"
	"github.com/italolelis/watchops/internal/app/provider/gh"
	"github.com/italolelis/watchops/internal/app/provider/gitlab"
	"github.com/italolelis/watchops/internal/app/provider/opsgenie"
	"github.com/italolelis/watchops/internal/app/provider/pagerduty"
	"github.com/italolelis/watchops/internal/app/storage"
	"github.com/italolelis/watchops/internal/app/stream"
	"github.com/italolelis/watchops/internal/app/subscriber"
	"github.com/italolelis/watchops/internal/pkg/log"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	LogLevel string `split_words:"true" default:"info"`
	Database struct {
		Driver          string `default:"postgres"`
		DSN             string
		MaxOpenConns    int           `split_words:"true" default:"30"`
		MaxIdleConns    int           `split_words:"true" default:"5"`
		ConnMaxLifetime time.Duration `split_words:"true" default:"1h"`
		Timeout         time.Duration `required:"true" default:"30s"`
		SchemaName      string        `split_words:"true" default:"watchops"`

		Bigquery struct {
			ProjectID string `split_words:"true"`
		}
	}
	MessageBroker subscriber.Config `split_words:"true"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := log.WithContext(context.Background())
	defer logger.Sync()

	if err := run(ctx); err != nil {
		logger.Fatalw("application finished with an error", "error", err)
	}
}

func run(ctx context.Context) error {
	logger := log.WithContext(ctx)
	// =============================================
	// Load Configuration
	// =============================================
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		return fmt.Errorf("failed to load env var: %w", err)
	}

	log.SetLevel(cfg.LogLevel)

	// =========================================================================
	// Setup Databases
	// =========================================================================
	logger.Debugw("connecting to event store", "driver", cfg.Database.Driver)
	db, err := storage.Connect(ctx, storage.Config(cfg.Database))
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	defer db.Close()

	logger.Debugw("creating subscriber", "driver", cfg.MessageBroker.Driver)
	subs, err, shutdown := subscriber.Build(ctx, cfg.MessageBroker.Driver, subscriber.Config(cfg.MessageBroker))
	if err != nil {
		return fmt.Errorf("failed to build subscriber: %w", err)
	}
	defer shutdown()

	pr := provider.NewParserRegistry()
	pr.
		Register(&gh.Parser{}).
		Register(&opsgenie.Parser{}).
		Register(&gh.Parser{}).
		Register(&gitlab.Parser{}).
		Register(&pagerduty.Parser{})

	events := stream.NewEventDataHandler(db, pr)
	return subs.Subscribe(ctx, events.Handle)
}
