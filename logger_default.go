package logger

import (
	"context"
	"encoding/json"

	"github.com/spf13/cast"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type defaultLogger struct {
	zapLogger *zap.Logger
}

func NewLogger(opt Option) Logger {
	return &defaultLogger{
		zapLogger: NewZapLogger(opt),
	}
}

func (d *defaultLogger) Debug(ctx context.Context, message string, details ...interface{}) {
	zapLogs := []zap.Field{
		zap.String("level", "debug"),
	}

	traceContextFields := TraceContext(ctx)

	fields := formatToField(details...)
	zapLogs = append(zapLogs, formatLogs(ctx, fields...)...)
	d.zapLogger.With(traceContextFields...).Debug(message, zapLogs...)
}

func (d *defaultLogger) Info(ctx context.Context, message string, details ...interface{}) {
	zapLogs := []zap.Field{}

	traceContextFields := TraceContext(ctx)

	fields := formatToField(details...)
	zapLogs = append(zapLogs, formatLogs(ctx, fields...)...)
	d.zapLogger.With(traceContextFields...).Info(message, zapLogs...)
}

func (d *defaultLogger) Warn(ctx context.Context, message string, details ...interface{}) {
	zapLogs := []zap.Field{}

	traceContextFields := TraceContext(ctx)

	fields := formatToField(details...)
	zapLogs = append(zapLogs, formatLogs(ctx, fields...)...)
	d.zapLogger.With(traceContextFields...).Warn(message, zapLogs...)
}

func (d *defaultLogger) Error(ctx context.Context, message string, details ...interface{}) {
	zapLogs := []zap.Field{}

	traceContextFields := TraceContext(ctx)

	fields := formatToField(details...)
	zapLogs = append(zapLogs, formatLogs(ctx, fields...)...)
	d.zapLogger.With(traceContextFields...).Error(message, zapLogs...)
}

func (d *defaultLogger) Fatal(ctx context.Context, message string, details ...interface{}) {
	zapLogs := []zap.Field{}

	traceContextFields := TraceContext(ctx)

	fields := formatToField(details...)
	zapLogs = append(zapLogs, formatLogs(ctx, fields...)...)
	d.zapLogger.With(traceContextFields...).Fatal(message, zapLogs...)
}

func (d *defaultLogger) Panic(ctx context.Context, message string, details ...interface{}) {
	zapLogs := []zap.Field{}

	traceContextFields := TraceContext(ctx)

	fields := formatToField(details...)
	zapLogs = append(zapLogs, formatLogs(ctx, fields...)...)
	d.zapLogger.With(traceContextFields...).Panic(message, zapLogs...)
}

func formatToField(details ...interface{}) (logRecord []Field) {
	for index, msg := range details {
		logRecord = append(logRecord, Field{
			Key: "message_" + cast.ToString(index),
			Val: msg,
		})
	}

	return
}

func formatLogs(ctx context.Context, fields ...Field) (logRecord []zap.Field) {
	ctxVal := ExtractCtx(ctx)

	// Add global value from context that must be exist on all logs!
	logRecord = append(logRecord, zap.String("app_name", ctxVal.ServiceName))
	logRecord = append(logRecord, zap.String("app_version", ctxVal.ServiceVersion))
	logRecord = append(logRecord, zap.Int("app_port", ctxVal.ServicePort))
	logRecord = append(logRecord, zap.String("app_tag", ctxVal.Tag))
	logRecord = append(logRecord, zap.String("app_method", ctxVal.ReqMethod))
	logRecord = append(logRecord, zap.String("app_uri", ctxVal.ReqURI))

	// Add additional data that available across all log, such as user_id
	if ctxVal.AdditionalData != nil {
		logRecord = append(logRecord, zap.Any("app_data", ctxVal.AdditionalData))
	}

	for _, field := range fields {
		logRecord = append(logRecord, formatLog(field.Key, field.Val))
	}

	return
}

func formatLog(key string, msg interface{}) (logRecord zap.Field) {
	if msg == nil {
		logRecord = zap.Any(key, struct{}{})
		return
	}

	// Try to convert the proto message into json
	p, ok := msg.(proto.Message)
	if ok {
		b, _ := json.Marshal(p)

		var data interface{}
		if err := json.Unmarshal(b, &data); err != nil {
			// Fallback: Just print the message
			// The message format should be like:
			// "name:\"mamatosai\"  id:1234  email:\"Rachmat Adi Prakoso\""
			logRecord = zap.Any(key, p)
			return
		}

		logRecord = zap.Any(key, data)
		return
	}

	// Try to convert the proto message into json
	// Incase the message is json stringify
	str, ok := msg.(string)
	if ok {
		var data interface{}
		if err := json.Unmarshal([]byte(str), &data); err != nil {
			// Fallback: Just print the original message
			logRecord = zap.Any(key, str)
			return
		}

		logRecord = zap.Any(key, data)
		return
	}

	logRecord = zap.Any(key, msg)
	return
}
