package gotel

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func NewOtelWithJaegerExporter(serviceName string, opt OtelWithJaegerOption) Gotel {
	ctx := context.Background()

	traceOpts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(opt.Endpoint),
	}

	if !opt.IsSecure {
		traceOpts = append(traceOpts, otlptracehttp.WithInsecure())
	}

	exporter, err := otlptracehttp.New(ctx, traceOpts...)
	if err != nil {
		panic(err)
	}

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(serviceName),
			),
		),
	)

	otel.SetTracerProvider(traceProvider)

	// Set the propagator
	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	gotel := &gotel{
		tracer:         traceProvider.Tracer(serviceName),
		tracerProvider: traceProvider,
		propagator:     prop,
	}

	// Set to global for easy to use every where
	_gotel = gotel

	return gotel
}
