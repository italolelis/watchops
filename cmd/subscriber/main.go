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

	events := stream.NewEventDataHandler(db, &gh.Parser{})
	subs, err := subscriber.Build(ctx, "kinesis", map[string]string{
		"region":   "eu-central-1",
		"endpoint": "http://localhost:4566",
	})

	return subs.Subscribe(ctx, "fourkeys_github", events.Handle)
}
