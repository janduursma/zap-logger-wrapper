# zap-logger-wrapper

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/janduursma/zap-logger-wrapper)](https://goreportcard.com/report/github.com/janduursma/zap-logger-wrapper)
[![semantic-release](https://img.shields.io/badge/semantic--release-ready-brightgreen)](https://github.com/go-semantic-release/go-semantic-release)
[![codecov](https://codecov.io/gh/janduursma/zap-logger-wrapper/graph/badge.svg?token=NA24AOA3EN)](https://codecov.io/gh/janduursma/zap-logger-wrapper)

This package provides a convenience wrapper around [Uber Zap](https://github.com/uber-go/zap) to streamline:

- **Logger initialization** (with production-friendly defaults and a “service” field).
- **Automatic injection of a Trace ID** (via a user-supplied function).
- **Structured, leveled logging** (Info, Error, Debug, etc.).
- **Optional** sub-logger creation with extra default fields (`With` method).

---

## Features

1. **Easy Setup**  
   One call to `New(...)` returns a ready-to-use logger that writes JSON logs to `stdout` (or any configured output path).

2. **Trace ID Injection**  
   Supply a `GetTraceIDFn` that knows how to extract a trace ID from `context.Context`, so logs automatically include `"trace_id"`.

3. **Structured Fields**  
   The wrapper uses Zap’s _SugaredLogger_, so you can quickly add key-value pairs to each log call.

4. **Sub-Loggers**  
   Use `With(...)` to create loggers that always include certain fields (e.g., `"component":"myModule"`).

---

## Installation

```bash
go get github.com/janduursma/zap-logger-wrapper/v2
```

Then import it:

```go
import logger "github.com/janduursma/zap-logger-wrapper/v2"
```

---

## Configuration

This logger wrapper uses functional options to allow you to customize its behavior. By default, the logger is configured as follows:

- **Log Level:** Info  
  The default log level is set to `Info`. You can override this using the `WithLogLevel` option.

- **Output Paths:** `["stdout"]`  
  The default output path is set to standard output. Use `WithOutputPaths` to direct logs to a file or other destinations.

- **GetTraceIDFn:** `nil`  
  By default, no trace ID is automatically added to logs. If you want to include trace IDs (for example, when using distributed tracing), use `WithGetTraceIDFn` to supply a custom function that extracts the trace ID from your context.

---

## Usage

Below is a minimal usage example:

```go
package main

import (
   "context"
   "fmt"
   "os"

   logger "github.com/janduursma/zap-logger-wrapper/v2"
   "go.uber.org/zap"
)

func main() {
   ctx := context.Background()

   // 1) Optionally define a function that extracts trace IDs from a context.
   traceFn := func(ctx context.Context) string {
      return "fake-trace-id" // simple placeholder
   }

   // 2) Create a logger.
   log, err := logger.New("myServiceName", logger.WithTraceID(traceFn))
   if err != nil {
      log.Error(ctx, "failed to create logger", "err", err)
      os.Exit(1)
   }

   // 3) Log some messages with key-value pairs.
   log.Info(ctx, "Hello from logger", "key", "value")
   log.Error(ctx, "An error occurred", "code", 500)

   // 4) Optionally create a sub-logger with extra default fields.
   subLogger := log.With("component", "signup")
   subLogger.Debug(ctx, "Debug details here", "anotherField", true)

   // 5) Flush logs before exiting
   if err := log.Sync(); err != nil {
      fmt.Printf("failed to sync logs: %v\n", err)
   }
}
```

---

## Output

When run, you'll see structured logs (JSON) on stdout, like:

```json
{"level":"info","ts":"2025-01-01T12:34:56.789Z","msg":"Hello from logger","key":"value","trace_id":"fake-trace-id","service":"myServiceName"}
{"level":"error","ts":"2025-01-01T12:34:56.790Z","msg":"An error occurred","code":500,"trace_id":"fake-trace-id","service":"myServiceName"}
{"level":"debug","ts":"2025-01-01T12:34:56.791Z","msg":"Debug details here","anotherField":true,"component":"signup","trace_id":"fake-trace-id","service":"myServiceName"}
```

---

## Running Tests
```sh
go test ./...
```

---

## License
- [MIT License](LICENSE)
