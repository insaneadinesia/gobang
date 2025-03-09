package main

import (
	"context"
	"time"

	"github.com/insaneadinesia/gobang/gotel"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	// Configure Jaeger exporter options
	jaegerOpt := gotel.OtelWithJaegerOption{
		Endpoint: "localhost:4318",
		IsSecure: false,
	}

	// Initialize tracer provider
	tp := gotel.NewOtelWithJaegerExporter("inventory-service", jaegerOpt)
	if provider, ok := tp.DefaultTracerProvider().(*sdktrace.TracerProvider); ok {
		defer provider.Shutdown(context.Background())
	}

	tracer := tp.DefaultTracer()

	// Create a root span
	ctx, mainSpan := tracer.Start(context.Background(), "process-order")
	defer mainSpan.End()

	// Simulate work with child spans
	for i := range 3 {
		ctx, childSpan := tracer.Start(ctx, "item-processing")
		childSpan.SetAttributes(
			attribute.Int("item.index", i),
			attribute.String("worker.id", "test-worker-id"),
		)

		// Simulate processing time
		time.Sleep(50 * time.Millisecond)

		if i%2 == 0 {
			// Create nested span
			_, specialSpan := tracer.Start(ctx, "special-handling")
			time.Sleep(20 * time.Millisecond)
			specialSpan.End()
		}

		childSpan.End()
	}
}
