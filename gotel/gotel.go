package gotel

import (
	"context"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type Gotel interface {
	DefaultTracer() trace.Tracer
	DefaultTracerProvider() trace.TracerProvider
	ExtractCarier(carrier propagation.MapCarrier) context.Context
	InjectCarier(ctx context.Context, carrier propagation.MapCarrier)
	GetTextMapPropagator() propagation.TextMapPropagator
}

type gotel struct {
	tracer         trace.Tracer
	tracerProvider trace.TracerProvider
	propagator     propagation.TextMapPropagator
}

func (g *gotel) DefaultTracer() trace.Tracer {
	return g.tracer
}

func (g *gotel) DefaultTracerProvider() trace.TracerProvider {
	return g.tracerProvider
}

func (g *gotel) ExtractCarier(carrier propagation.MapCarrier) context.Context {
	return g.propagator.Extract(context.Background(), carrier)
}

func (g *gotel) InjectCarier(ctx context.Context, carrier propagation.MapCarrier) {
	g.propagator.Inject(ctx, carrier)
}

func (g *gotel) GetTextMapPropagator() propagation.TextMapPropagator {
	return g.propagator
}
