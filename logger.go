// Package logger provides a wrapper around Uber's Zap logger
// with automatic trace ID injection and structured logging.
package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// GetTraceIDFn is a function type that, given a context.Context, returns a trace ID.
type GetTraceIDFn func(ctx context.Context) string

// Logger is the wrapper around zap.SugaredLogger.
type Logger struct {
	zapLogger  *zap.SugaredLogger
	getTraceID GetTraceIDFn
}

// New creates a new Logger wrapper. The arguments are:
//   - zapLogger: the underlying zap.Logger
//   - fn: a function that returns a trace ID from the given context
func New(service string, fn GetTraceIDFn, level zapcore.Level, outputPaths ...string) (*Logger, error) {
	config := zap.NewProductionConfig()

	config.Level = zap.NewAtomicLevelAt(level)
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]any{
		"service": service,
	}

	config.OutputPaths = []string{"stdout"}
	if outputPaths != nil {
		config.OutputPaths = outputPaths
	}

	log, err := config.Build(zap.WithCaller(true))
	if err != nil {
		return nil, err
	}

	return &Logger{
		zapLogger:  log.Sugar(),
		getTraceID: fn,
	}, nil
}

// Info logs a message at InfoLevel, automatically including trace_id if available.
func (l *Logger) Info(ctx context.Context, msg string, keyVals ...interface{}) {
	if l.getTraceID != nil {
		if traceID := l.getTraceID(ctx); traceID != "" {
			// Append the trace_id as a key-value pair
			keyVals = append(keyVals, "trace_id", traceID)
		}
	}
	l.zapLogger.Infow(msg, keyVals...)
}

// Error logs a message at ErrorLevel, automatically including trace_id if available.
func (l *Logger) Error(ctx context.Context, msg string, keyVals ...interface{}) {
	if l.getTraceID != nil {
		if traceID := l.getTraceID(ctx); traceID != "" {
			keyVals = append(keyVals, "trace_id", traceID)
		}
	}
	l.zapLogger.Errorw(msg, keyVals...)
}

// Debug logs a message at DebugLevel, automatically including trace_id if available.
func (l *Logger) Debug(ctx context.Context, msg string, keyVals ...interface{}) {
	if l.getTraceID != nil {
		if traceID := l.getTraceID(ctx); traceID != "" {
			keyVals = append(keyVals, "trace_id", traceID)
		}
	}
	l.zapLogger.Debugw(msg, keyVals...)
}

// With returns a child Logger that includes some default key-value pairs.
// The child still auto-injects trace IDs.
func (l *Logger) With(keyVals ...interface{}) *Logger {
	// zap.SugaredLogger has a With(...) method that returns a new SugaredLogger
	newSugared := l.zapLogger.With(keyVals...)
	return &Logger{
		zapLogger:  newSugared,
		getTraceID: l.getTraceID,
	}
}

// Sync flushes any buffered log entries.
func (l *Logger) Sync() error {
	return l.zapLogger.Sync()
}
