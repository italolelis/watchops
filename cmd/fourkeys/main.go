//nolint: exhaustivestruct
package main

import (
	"context"
	"expvar"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/heptiolabs/healthcheck"
	"github.com/italolelis/fourkeys/internal/app/http/rest"
	"github.com/italolelis/fourkeys/internal/app/provider/gh"
	"github.com/italolelis/fourkeys/internal/app/provider/opsgenie"
	"github.com/italolelis/fourkeys/internal/app/publisher"
	"github.com/italolelis/fourkeys/internal/app/publisher/kinesis"
	"github.com/italolelis/fourkeys/internal/app/wh"
	"github.com/italolelis/fourkeys/internal/pkg/log"
	"github.com/italolelis/fourkeys/internal/pkg/signal"
	"github.com/kelseyhightower/envconfig"
)

var version = "develop"

const goRoutineCount = 100

type config struct {
	LogLevel string `split_words:"true"`
	Web      struct {
		APIHost         string        `split_words:"true" default:"0.0.0.0:8080"`
		ProbeHost       string        `split_words:"true" default:"0.0.0.0:9090"`
		ReadTimeout     time.Duration `split_words:"true" default:"30s"`
		WriteTimeout    time.Duration `split_words:"true" default:"30s"`
		IdleTimeout     time.Duration `split_words:"true" default:"5s"`
		ShutdownTimeout time.Duration `split_words:"true" default:"30s"`
	}
	Github struct {
		WebhookSecret string `split_words:"true"`
	}
	Opsgenie struct {
		WebhookSecret string `split_words:"true"`
	}
	Kinesis struct {
		Endpoint   string
		Region     string `default:"eu-central-1"`
		DisableSSL bool   `split_words:"true"`
	}
	TopicPrefix string `split_words:"true"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := log.WithContext(ctx)
	defer logger.Sync()

	if err := run(ctx); err != nil {
		logger.Fatal(err)
	}
}

func run(ctx context.Context) error {
	// =========================================================================
	// Logging
	var (
		logger = log.WithContext(ctx).Named("main")
		slog   = logger.Named("starting")
	)

	// =========================================================================
	// Configuration
	var cfg config

	if err := envconfig.Process("", &cfg); err != nil {
		return fmt.Errorf("failed to load the env vars: %w", err)
	}

	log.SetLevel(cfg.LogLevel)

	// =============================================
	// Kinesis
	// =============================================
	slog.Info("kinesis")

	session, err := kinesis.NewSession(kinesis.SessionConfig{
		Endpoint:   cfg.Kinesis.Endpoint,
		Region:     cfg.Kinesis.Region,
		DisableSSL: cfg.Kinesis.DisableSSL,
	})
	if err != nil {
		return fmt.Errorf("failed to create a new aws session: %w", err)
	}

	p, err := kinesis.NewPublisher(session)
	if err != nil {
		return fmt.Errorf("failed to create a new publisher: %w", err)
	}

	// =========================================================================
	// App Starting
	expvar.NewString("build").Set(version)
	slog.Infow("Application initializing", "version", version)

	defer logger.Info("completed")

	// =========================================================================
	// Start API Service

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	server := setupServer(ctx, p, cfg)

	go func() {
		slog.Infow("Initializing API support", "host", cfg.Web.APIHost)
		serverErrors <- server.ListenAndServe()
	}()

	// =========================================================================
	// Start Observability API Service
	probeServer := setupProbeServer(cfg)

	go func() {
		slog.Infow("Initializing probe API support", "host", cfg.Web.ProbeHost)
		serverErrors <- probeServer.ListenAndServe()
	}()

	done := signal.New(ctx)

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case <-done.Done():
		logger.Infow("start shutdown")

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Errorw("failed to gracefully shutdown the server", "err", err)

			if err = server.Close(); err != nil {
				return fmt.Errorf("could not stop server gracefully: %w", err)
			}
		}

		if err := probeServer.Shutdown(ctx); err != nil {
			logger.Errorw("failed to gracefully shutdown the probe server", "err", err)

			if err = probeServer.Close(); err != nil {
				return fmt.Errorf("could not stop probe server gracefully: %w", err)
			}
		}
	}

	return nil
}

func setupProbeServer(cfg config) *http.Server {
	h := healthcheck.NewHandler()
	// Our app is not happy if we've got more than 100 goroutines running.
	h.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(goRoutineCount))

	return &http.Server{
		Addr:    cfg.Web.ProbeHost,
		Handler: h,
	}
}

// setupServer prepares the handlers and services to create the http rest server.
func setupServer(ctx context.Context, p publisher.Publisher, cfg config) *http.Server {
	webhook := wh.NewPublisherConnector(ctx, p, cfg.TopicPrefix)
	ghHandler := rest.NewWebhookHandler(webhook)

	ghValidator := gh.NewValidator(cfg.Github.WebhookSecret)
	ogValidator := opsgenie.NewValidator(cfg.Opsgenie.WebhookSecret)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/webhooks", ghHandler.Routes(ghValidator, ogValidator))

	return &http.Server{
		Addr:         cfg.Web.APIHost,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		Handler:      r,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}
}
