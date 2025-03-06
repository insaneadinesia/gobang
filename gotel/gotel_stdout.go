package gotel

import (
	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func NewOtelWithStdoutExporter(serviceName string) Gotel {
	exporter, err := stdout.New(stdout.WithPrettyPrint())
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

	gotel := &gotel{
		tracer:         traceProvider.Tracer(serviceName),
		tracerProvider: traceProvider,
	}

	// Set global variable
	Otel = gotel

	return gotel
}
