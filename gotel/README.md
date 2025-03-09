# Gotel - OpenTelemetry Integration for Go

A lightweight OpenTelemetry wrapper providing streamlined configuration for common tracing setups.

## Features

- Simplified tracer provider configuration
- Built-in exporters:
  - Stdout (development-friendly)
  - Jaeger (production-ready)
- Environment-based auto-configuration
- Type-safe configuration options

## Installation

```bash
go get github.com/insaneadinesia/gobang/gotel
```

## Quick Start

```go
package main

import (
    "context"
    "github.com/insaneadinesia/gobang/gotel"
)

func main() {
    // Initialize with stdout exporter
    tp := gotel.NewOtelWithStdoutExporter("example-service")
    defer tp.Shutdown(context.Background())

    tracer := tp.DefaultTracer()
    
    ctx, span := tracer.Start(context.Background(), "demo-span")
    defer span.End()
    
    // Your application logic here
}
```

## Configuration Options

```go
// Service name is set when creating the exporter:
gotel.NewOtelWithStdoutExporter("my-service")

// Jaeger configuration options:
gotel.OtelWithJaegerOption{
    Endpoint: "http://localhost:14268/api/traces", // Jaeger collector endpoint
    IsSecure: false, // Set true for HTTPS/TLS connections
}
```

## Advanced Example

See the [examples directory](/examples/advanced/main.go) for a complete demonstration with:
- Jaeger exporter configuration
- Custom sampling rates
- Environment variable integration
- Multiple span relationships

## Contributing

1. Fork the repository
2. Create a feature branch
3. Submit a PR with tests and documentation