package logger

import (
	"context"
	"strings"

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

func IsSkipPrintLog(msg string) bool {
	skipedContentType := []string{
		"Content-Type: application/",
		"Content-Type: audio/",
		"Content-Type: image/",
		"Content-Type: video/",
	}

	for _, contentType := range skipedContentType {
		if strings.Contains(msg, contentType) {
			return true
		}
	}

	return false
}
