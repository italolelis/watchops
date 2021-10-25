package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/italolelis/fourkeys/internal/app/provider/gh"
	"github.com/italolelis/fourkeys/internal/app/storage"
	"github.com/italolelis/fourkeys/internal/app/stream"
	"github.com/italolelis/fourkeys/internal/app/wh"
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
	lambda.Start(lambdaHandler(events.Handle))

	return nil
}

func lambdaHandler(fn stream.HandlerFunc) func(context.Context, events.KinesisEvent) error {
	return func(ctx context.Context, kinesisEvent events.KinesisEvent) error {
		logger := log.WithContext(ctx).Named("lambda")
		logger.Info("processing messages...")

		for _, record := range kinesisEvent.Records {
			logger.Debug("incoming message received")

			data := record.Kinesis.Data

			var messageContainer wh.MessageContainer
			if err := json.Unmarshal(data, &messageContainer); err != nil {
				return fmt.Errorf("failed to parse message from stream: %w", err)
			}

			messageContainer.Headers["msg_id"] = append(messageContainer.Headers["msg_id"], record.EventID)

			if err := fn(ctx, messageContainer.Payload, messageContainer.Headers); err != nil {
				logger.Errorw("failed to process event", "err", err)

				continue
			}

			logger.Debug("finished processing incoming message")
		}

		logger.Info("finished processing messages")

		return nil
	}
}
