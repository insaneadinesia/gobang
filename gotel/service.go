package gotel

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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

func ExtractCarier(carrier propagation.MapCarrier) context.Context {
	if _gotel == nil {
		return otel.GetTextMapPropagator().Extract(context.Background(), carrier)
	}

	return _gotel.ExtractCarier(carrier)
}

func InjectCarier(ctx context.Context, carrier propagation.MapCarrier) {
	if _gotel == nil {
		otel.GetTextMapPropagator().Inject(ctx, carrier)
		return
	}

	_gotel.InjectCarier(ctx, carrier)
}

func GetTextMapPropagator() propagation.TextMapPropagator {
	if _gotel == nil {
		return otel.GetTextMapPropagator()
	}
	return _gotel.GetTextMapPropagator()
}
