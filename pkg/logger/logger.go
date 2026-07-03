package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"

	"real-time-forum/pkg/config"
)

// stackTraceHandler wraps an slog.Handler and automatically attaches
// a stack trace to ERROR-level log entries.
type stackTraceHandler struct {
	handler slog.Handler
}

func (h *stackTraceHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *stackTraceHandler) Handle(ctx context.Context, rec slog.Record) error {
	if rec.Level == slog.LevelError {
		stackBuf := make([]byte, 4096)
		n := runtime.Stack(stackBuf, false)
		rec.AddAttrs(slog.String("stack", string(stackBuf[:n])))
	}
	return h.handler.Handle(ctx, rec)
}

func (h *stackTraceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &stackTraceHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *stackTraceHandler) WithGroup(name string) slog.Handler {
	return &stackTraceHandler{handler: h.handler.WithGroup(name)}
}

// InitLogger initializes a slog.Logger based on the app configuration.
// In production mode, it outputs JSON. In development, it outputs human-readable text.
// ERROR-level logs include an automatic stack trace.
func InitLogger(app *config.AppConfig) {
	level := slog.LevelInfo
	if app.LogLevel != "" {
		switch app.LogLevel {
		case "debug":
			level = slog.LevelDebug
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		}
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: !app.InProduction,
	}

	var baseHandler slog.Handler
	if app.InProduction {
		baseHandler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		baseHandler = slog.NewTextHandler(os.Stdout, opts)
	}

	app.Logger = slog.New(&stackTraceHandler{handler: baseHandler})
	slog.SetDefault(app.Logger)
}