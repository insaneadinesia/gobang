package logger

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TraceContext(ctx context.Context) []zapcore.Field {
	const (
		FieldKeyTraceID = "trace.id"
		FieldKeySpanID  = "span.id"
	)

	spanContext := trace.SpanFromContext(ctx).SpanContext()
	fields := []zapcore.Field{
		zap.Stringer(FieldKeyTraceID, spanContext.TraceID()),
		zap.Stringer(FieldKeySpanID, spanContext.SpanID()),
	}

	return fields
}
