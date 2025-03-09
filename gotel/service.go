package gotel

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var _gotel Gotel

func DefaultTracer() trace.Tracer {
	if _gotel == nil {
		return otel.GetTracerProvider().Tracer("")
	}

	return _gotel.DefaultTracer()
}

func DefaultTracerProvider() trace.TracerProvider {
	if _gotel == nil {
		return otel.GetTracerProvider()
	}

	return _gotel.DefaultTracerProvider()
}
