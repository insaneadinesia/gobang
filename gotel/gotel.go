package gotel

import (
	"go.opentelemetry.io/otel/trace"
)

type Gotel interface {
	DefaultTracer() trace.Tracer
	DefaultTracerProvider() trace.TracerProvider
}

type gotel struct {
	tracer         trace.Tracer
	tracerProvider trace.TracerProvider
}

func (g *gotel) DefaultTracer() trace.Tracer {
	return g.tracer
}

func (g *gotel) DefaultTracerProvider() trace.TracerProvider {
	return g.tracerProvider
}
