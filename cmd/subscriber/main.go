package main

import (
	"context"
	"fmt"
	"time"

	"github.com/italolelis/fourkeys/internal/app/provider/gh"
	"github.com/italolelis/fourkeys/internal/app/storage"
	"github.com/italolelis/fourkeys/internal/app/stream"
	"github.com/italolelis/fourkeys/internal/app/subscriber"
	"github.com/italolelis/fourkeys/internal/pkg/log"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	LogLevel string `split_words:"true"`
	Database struct {
		Driver          string        `default:"postgres"`
		DSN             string        `required:"true"`
		MaxOpenConns    int           `split_words:"true" default:"30"`
		MaxIdleConns    int           `split_words:"true" default:"5"`
		ConnMaxLifetime time.Duration `split_words:"true" default:"1h"`
		Timeout         time.Duration `required:"true" default:"30s"`
	}
	MessageBroker struct {
		Driver  string `required:"true"`
		Kinesis struct {
			Endpoint   string
			Region     string `required:"true"`
			StreamName string `split_words:"true" required:"true"`
		}
	} `split_words:"true"`
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
	db, err := storage.Connect(ctx, storage.Config(cfg.Database))
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	defer db.Close()

	logger := log.WithContext(context.Background())
	logger.Debugw("selecting subscriber driver", "driver", cfg.MessageBroker.Driver)
	subs, err := subscriber.Build(ctx, "kinesis", subscriber.Config(cfg.MessageBroker))

	events := stream.NewEventDataHandler(db, &gh.Parser{})
	return subs.Subscribe(ctx, "fourkeys_github", events.Handle)
}
