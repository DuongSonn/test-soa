package logger

import (
	"context"
	"log/slog"
	"os"
)

var defaultLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	AddSource: true,
	ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == "password" {
			a.Value = slog.StringValue("***")
		}
		return a
	},
}))

// contextKey is used to store and retrieve values from context
type contextKey string

const (
	requestIDKey contextKey = "request_id"
	userIDKey    contextKey = "user_id"
)

// Logger wraps slog.Logger with context
type Logger struct {
	ctx context.Context
	log *slog.Logger
}

// WithCtx creates a new Logger with context
func WithCtx(ctx context.Context) *Logger {
	return &Logger{ctx: ctx, log: defaultLogger}
}

// WithRequestID adds request ID to context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// WithUserID adds user ID to context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// getContextAttrs extracts logging attributes from context
func (l *Logger) getContextAttrs() []any {
	var attrs []any
	if l.ctx == nil {
		return attrs
	}

	if reqID := l.ctx.Value(requestIDKey); reqID != nil {
		attrs = append(attrs, slog.String("request_id", reqID.(string)))
	}
	if userID := l.ctx.Value(userIDKey); userID != nil {
		attrs = append(attrs, slog.String("user_id", userID.(string)))
	}
	return attrs
}

// Info logs at INFO level
func (l *Logger) Info(msg string, attrs ...any) {
	ctxAttrs := l.getContextAttrs()
	l.log.Info(msg, append(ctxAttrs, attrs...)...)
}

// Error logs at ERROR level
func (l *Logger) Error(msg string, attrs ...any) {
	ctxAttrs := l.getContextAttrs()
	l.log.Error(msg, append(ctxAttrs, attrs...)...)
}

// Debug logs at DEBUG level
func (l *Logger) Debug(msg string, attrs ...any) {
	ctxAttrs := l.getContextAttrs()
	l.log.Debug(msg, append(ctxAttrs, attrs...)...)
}

// Warn logs at WARN level
func (l *Logger) Warn(msg string, attrs ...any) {
	ctxAttrs := l.getContextAttrs()
	l.log.Warn(msg, append(ctxAttrs, attrs...)...)
}
