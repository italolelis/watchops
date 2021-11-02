package log

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

type traceIDKeyType int

const (
	traceIDKey    traceIDKeyType = iota
	TraceIDHeader                = "X-B3-TraceId"
)

func NewStructuredLogger(logger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{logger})
}

type StructuredLogger struct {
	Logger *zap.SugaredLogger
}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	traceID := r.Header.Get(TraceIDHeader)

	// adds the trace id into the request context.
	ctx := context.WithValue(r.Context(), traceIDKey, traceID)
	r = r.WithContext(ctx)

	logger = logger.With(
		"trace_id", traceID,
		"method", r.Method,
		"host", r.Host,
		"request", r.RequestURI,
		"remote-addr", r.RemoteAddr,
		"referer", r.Referer(),
		"user-agent", r.UserAgent(),
	)

	logger.Info("request started")

	return &StructuredLoggerEntry{Logger: l.Logger}
}

type StructuredLoggerEntry struct {
	Logger *zap.SugaredLogger
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger = l.Logger.With("resp_status", status, "resp_bytes_length", bytes, "resp_elapsed_ms", float64(elapsed.Nanoseconds())/1000000.0)

	l.Logger.Info("request complete")
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.With(zap.Field{
		Key:    "stack",
		String: string(stack),
	}, zap.Field{
		Key:    "panic",
		String: fmt.Sprintf("%+v", v),
	})
}
