package logger

import (
	"context"

	iContext "github.com/eduardohoraciosanto/users/internal/context"
	"go.uber.org/zap"
)

type Logger interface {
	Info(ctx context.Context, message string)
	Error(ctx context.Context, message string)
	Debug(ctx context.Context, message string)
	Warn(ctx context.Context, message string)
	WithField(key string, value interface{}) Logger
	WithError(err error) Logger
	Sync() error
}

type logger struct {
	internal       *zap.SugaredLogger
	tracingEnabled bool
}

// NewLogger initializes de logger and allows the usage of the Logger Interface
// user MUST defer the call to Sync() immediately afterwards.

func NewLogger(service string, version string, tracingEnabled bool) Logger {
	l, err := zap.NewProduction()
	if err != nil {
		panic("unable to initialize logger: " + err.Error())
	}

	return &logger{
		internal: l.WithOptions(zap.AddCallerSkip(1)).Sugar().
			With("version", version).
			With("service", service),
		tracingEnabled: tracingEnabled,
	}
}

// Sync MUST be deferred to flush any buffered logs prior to shutting down the application.
func (l *logger) Sync() error {
	return l.internal.Sync()
}

// Info allows for a message with info lever to be logged
func (l *logger) Info(ctx context.Context, message string) {
	if l.tracingEnabled {
		l.injectTracing(ctx).Info(message)
		return
	}
	l.internal.Info(message)
}

// Error allows for a message with error lever to be logged
func (l *logger) Error(ctx context.Context, message string) {
	if l.tracingEnabled {
		l.injectTracing(ctx).Error(message)
		return
	}
	l.internal.Error(message)
}

// Debug allows for a message with debug lever to be logged
func (l *logger) Debug(ctx context.Context, message string) {
	if l.tracingEnabled {
		l.injectTracing(ctx).Debug(message)
		return
	}
	l.internal.Debug(message)
}

// Warn allows for a message with warn lever to be logged
func (l *logger) Warn(ctx context.Context, message string) {
	if l.tracingEnabled {
		l.injectTracing(ctx).Warn(message)
		return
	}
	l.internal.Warn(message)

}

// WithField allows for the inclusion of a key-value into the log
func (l *logger) WithField(key string, value interface{}) Logger {
	return &logger{
		internal:       l.internal.With(key, value),
		tracingEnabled: l.tracingEnabled,
	}
}

// WithError allows for the inclusion of an error into the log.
func (l *logger) WithError(err error) Logger {
	newLogger := l.internal.With(
		"error", err,
	)

	return &logger{
		internal:       newLogger,
		tracingEnabled: l.tracingEnabled,
	}
}

// injectTracing enters correlation ID if any.
func (l *logger) injectTracing(ctx context.Context) *zap.SugaredLogger {
	//add our correlation id if present
	cid := ctx.Value(iContext.CorrelationID("correlation_id"))
	rIP := ctx.Value(iContext.RemoteIP("remote_ip"))
	entry := l.internal
	if cid != nil {
		entry = entry.With("correlation_id", ctx.Value(iContext.CorrelationID("correlation_id")))
	}
	if rIP != nil {
		entry = entry.With("remote_ip", ctx.Value(iContext.RemoteIP("remote_ip")))
	}

	return entry
}
