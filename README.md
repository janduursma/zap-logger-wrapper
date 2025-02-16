# zap-logger-wrapper

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/janduursma/zap-logger-wrapper)](https://goreportcard.com/report/github.com/janduursma/zap-logger-wrapper)
[![semantic-release](https://img.shields.io/badge/semantic--release-ready-brightgreen)](https://github.com/go-semantic-release/go-semantic-release)
[![codecov](https://codecov.io/gh/janduursma/zapnote-logger-go/graph/badge.svg?token=NA24AOA3EN)](https://codecov.io/gh/janduursma/zapnote-logger-go)

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
go get github.com/janduursma/zap-logger-wrapper
```

Then import it:

```go
import logger "github.com/janduursma/zap-logger-wrapper"
```

---

## Usage

Below is a minimal usage example:

```go
package main

import (
   "context"
   "fmt"
   "log"
   
   logger "github.com/janduursma/zap-logger-wrapper"
   "go.uber.org/zap"
)

func main() {
   // 1) Define a function that extracts trace IDs from a context.
   traceFn := func(ctx context.Context) string {
      // e.g., if stored under a custom key:
      // if v, ok := ctx.Value("traceIDKey").(string); ok {
      //     return v
      // }
      // return ""
      return "fake-trace-id" // simple placeholder
   }

   // 2) Create a logger (by default logs go to stdout using JSON).
   l, err := logger.New("myServiceName", traceFn, zap.InfoLevel)
   if err != nil {
      log.Fatalf("failed to create logger: %v", err)
   }

   ctx := context.Background()
   // Optionally, store a real trace ID in context:
   // ctx = context.WithValue(ctx, "traceIDKey", "abcd1234")

   // 3) Log some messages with key-value pairs.
   l.Info(ctx, "Hello from logger", "key", "value")
   l.Error(ctx, "An error occurred", "code", 500)

   // 4) Optionally create a sub-logger with extra default fields.
   subLogger := l.With("component", "signup")
   subLogger.Debug(ctx, "Debug details here", "anotherField", true)

   // 5) Flush logs before exiting
   if err := l.Sync(); err != nil {
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
