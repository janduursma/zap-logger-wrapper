// Package logger provides a wrapper around Uber's Zap logger
// It supports configuring options such as the log level, output paths, and a GetTraceIDFn
// to automatically include trace IDs in log entries.
//
// By default, the logger is initialized with the following settings:
//   - Log Level: Info
//   - Output Paths: ["stdout"]
//   - GetTraceIDFn: nil (i.e. no trace ID is automatically added)
//
// These defaults can be overridden using the provided functional options.
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
	zapLogger    *zap.SugaredLogger
	getTraceIDFn GetTraceIDFn
	level        zapcore.Level
	outputPaths  []string
}

// Option defines a functional option for configuring the Logger.
type Option func(manager *Logger)

// WithTraceID allows a custom function to be set which automatically adds trace IDs to logs.
func WithTraceID(getTraceIDFn GetTraceIDFn) Option {
	return func(l *Logger) {
		l.getTraceIDFn = getTraceIDFn
	}
}

// WithLevel allows a custom minimum logging level to be set.
func WithLevel(level zapcore.Level) Option {
	return func(l *Logger) {
		l.level = level
	}
}

// WithOutputPaths allows a custom output path to be set.
func WithOutputPaths(outputPaths []string) Option {
	return func(l *Logger) {
		l.outputPaths = outputPaths
	}
}

// New creates a new Logger wrapper around zap.SugaredLogger.
func New(service string, opts ...Option) (*Logger, error) {
	defaultTraceIDFn := func(_ context.Context) string { return "" }
	defaultLevel := zap.InfoLevel
	defaultOutputPaths := []string{"stdout"}

	l, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	logger := &Logger{
		zapLogger:    l.Sugar(),
		getTraceIDFn: defaultTraceIDFn,
		level:        defaultLevel,
		outputPaths:  defaultOutputPaths,
	}

	for _, opt := range opts {
		opt(logger)
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(logger.level)
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]any{
		"service": service,
	}
	config.OutputPaths = logger.outputPaths

	l, err = config.Build(zap.WithCaller(true))
	if err != nil {
		return nil, err
	}
	logger.zapLogger = l.Sugar()

	return logger, nil
}

// Info logs a message at InfoLevel, automatically including trace_id if available.
func (l *Logger) Info(ctx context.Context, msg string, keyVals ...interface{}) {
	if l.getTraceIDFn != nil {
		if traceID := l.getTraceIDFn(ctx); traceID != "" {
			// Append the trace_id as a key-value pair
			keyVals = append(keyVals, "trace_id", traceID)
		}
	}
	l.zapLogger.Infow(msg, keyVals...)
}

// Error logs a message at ErrorLevel, automatically including trace_id if available.
func (l *Logger) Error(ctx context.Context, msg string, keyVals ...interface{}) {
	if l.getTraceIDFn != nil {
		if traceID := l.getTraceIDFn(ctx); traceID != "" {
			keyVals = append(keyVals, "trace_id", traceID)
		}
	}
	l.zapLogger.Errorw(msg, keyVals...)
}

// Debug logs a message at DebugLevel, automatically including trace_id if available.
func (l *Logger) Debug(ctx context.Context, msg string, keyVals ...interface{}) {
	if l.getTraceIDFn != nil {
		if traceID := l.getTraceIDFn(ctx); traceID != "" {
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
		zapLogger:    newSugared,
		getTraceIDFn: l.getTraceIDFn,
		level:        l.level,
		outputPaths:  l.outputPaths,
	}
}

// Sync flushes any buffered log entries.
func (l *Logger) Sync() error {
	return l.zapLogger.Sync()
}
