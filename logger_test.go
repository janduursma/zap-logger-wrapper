package logger_test

import (
	"context"
	"net/url"
	"strings"
	"testing"

	logger "github.com/janduursma/zap-logger-wrapper"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// memorySink is a simple zapcore.WriteSyncer that stores logs in a string builder.
type memorySink struct {
	logs strings.Builder
}

// Write implements io.Writer.
func (m *memorySink) Write(p []byte) (n int, err error) {
	return m.logs.Write(p)
}

// Sync is a no-op for memorySink.
func (m *memorySink) Sync() error {
	return nil
}

// Close is a no-op for memorySink.
func (m *memorySink) Close() error {
	return nil
}

func TestLogger(t *testing.T) {
	// Register a custom sink, so we can specify the output path.
	sink := &memorySink{}
	require.NoError(t, zap.RegisterSink("test", func(_ *url.URL) (zap.Sink, error) {
		// 'u' is the parsed URL for "test://whatever"
		return sink, nil
	}), "failed to register test sink")

	// Create a logger via the public New(...) function,
	// overriding the output path to use the in-memory sink.
	traceFn := func(_ context.Context) string { return "test-trace-id" }
	l, err := logger.New("test-service", traceFn, zap.DebugLevel, "test://whatever")
	require.NoError(t, err, "failed to create logger")
	require.NotNil(t, l, "logger should not be nil")

	ctx := context.Background()

	// INFO
	l.Info(ctx, "Info message", "userID", 1234)

	// ERROR
	l.Error(ctx, "Error message", "errCode", 500)

	// Sub-logger with extra fields
	sub := l.With("component", "signup-flow")
	sub.Debug(ctx, "Debugging message", "extra", true)

	// Substring checks to confirm that the fields appear in the JSON.
	logs := sink.logs.String()

	// From config.InitialFields in logger.New(), we expect "service":"test-service"
	require.Contains(t, logs, `"service":"test-service"`, "service field should be present")

	// All logs should have a trace_id from traceFn
	require.Contains(t, logs, `"trace_id":"test-trace-id"`, "trace ID should be present")

	// Check each log level is present
	require.Contains(t, logs, `"level":"info"`, "info level should appear")
	require.Contains(t, logs, `"level":"error"`, "error level should appear")
	require.Contains(t, logs, `"level":"debug"`, "debug level should appear")

	// Check each message is present
	require.Contains(t, logs, `"msg":"Info message"`, "info message should be present")
	require.Contains(t, logs, `"msg":"Error message"`, "error message should be present")
	require.Contains(t, logs, `"msg":"Debugging message"`, "debug message should be present")

	// Additional fields
	require.Contains(t, logs, `"userID":1234`, "extra info field not found")
	require.Contains(t, logs, `"errCode":500`, "error code field not found")
	require.Contains(t, logs, `"component":"signup-flow"`, "sub-logger field not found")

	// Ensure no errors from Sync
	err = l.Sync()
	require.NoError(t, err, "logger.Sync() should not return an error")
}
