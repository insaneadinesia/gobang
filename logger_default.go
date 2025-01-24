package logger

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/spf13/cast"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type defaultLogger struct {
	zapLogger     *zap.Logger
	enableMasking bool
	maskingFields map[string]bool
}

func NewLogger(opt Option) Logger {
	// Assign masking fields to global variable
	maskingFields := make(map[string]bool)
	for _, fieldName := range opt.MaskingFields {
		maskingFields[fieldName] = true
	}

	return &defaultLogger{
		zapLogger:     NewZapLogger(opt),
		enableMasking: opt.EnableMaskingFields,
		maskingFields: maskingFields,
	}
}

func (d *defaultLogger) Debug(ctx context.Context, message string, details ...interface{}) {
	zapLogs := []zap.Field{
		zap.String("level", "debug"),
	}

	traceContextFields := TraceContext(ctx)

	fields := d.formatToField(details...)
	zapLogs = append(zapLogs, d.formatLogs(ctx, fields...)...)
	d.zapLogger.With(traceContextFields...).Debug(message, zapLogs...)
}

func (d *defaultLogger) Info(ctx context.Context, message string, details ...interface{}) {
	zapLogs := []zap.Field{}

	traceContextFields := TraceContext(ctx)

	fields := d.formatToField(details...)
	zapLogs = append(zapLogs, d.formatLogs(ctx, fields...)...)
	d.zapLogger.With(traceContextFields...).Info(message, zapLogs...)
}

func (d *defaultLogger) Warn(ctx context.Context, message string, details ...interface{}) {
	zapLogs := []zap.Field{}

	traceContextFields := TraceContext(ctx)

	fields := d.formatToField(details...)
	zapLogs = append(zapLogs, d.formatLogs(ctx, fields...)...)
	d.zapLogger.With(traceContextFields...).Warn(message, zapLogs...)
}

func (d *defaultLogger) Error(ctx context.Context, message string, details ...interface{}) {
	zapLogs := []zap.Field{}

	traceContextFields := TraceContext(ctx)

	fields := d.formatToField(details...)
	zapLogs = append(zapLogs, d.formatLogs(ctx, fields...)...)
	d.zapLogger.With(traceContextFields...).Error(message, zapLogs...)
}

func (d *defaultLogger) Fatal(ctx context.Context, message string, details ...interface{}) {
	zapLogs := []zap.Field{}

	traceContextFields := TraceContext(ctx)

	fields := d.formatToField(details...)
	zapLogs = append(zapLogs, d.formatLogs(ctx, fields...)...)
	d.zapLogger.With(traceContextFields...).Fatal(message, zapLogs...)
}

func (d *defaultLogger) Panic(ctx context.Context, message string, details ...interface{}) {
	zapLogs := []zap.Field{}

	traceContextFields := TraceContext(ctx)

	fields := d.formatToField(details...)
	zapLogs = append(zapLogs, d.formatLogs(ctx, fields...)...)
	d.zapLogger.With(traceContextFields...).Panic(message, zapLogs...)
}

func (d *defaultLogger) formatToField(details ...interface{}) (logRecord []Field) {
	for index, msg := range details {
		logRecord = append(logRecord, Field{
			Key: "message_" + cast.ToString(index),
			Val: msg,
		})
	}

	return
}

func (d *defaultLogger) formatLogs(ctx context.Context, fields ...Field) (logRecord []zap.Field) {
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
		logRecord = append(logRecord, d.formatLog(field.Key, field.Val))
	}

	return
}

func (d *defaultLogger) formatLog(key string, msg interface{}) (logRecord zap.Field) {
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

		logRecord = zap.Any(key, d.maskData(data))
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

		logRecord = zap.Any(key, d.maskData(data))
		return
	}

	logRecord = zap.Any(key, d.maskData(msg))
	return
}

func (d *defaultLogger) maskData(input interface{}) interface{} {
	if !d.enableMasking {
		return input
	}

	val := reflect.ValueOf(input)
	maskingFields := d.maskingFields

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Map:
		// If it's a map, iterate through the keys and mask sensitive ones
		for _, key := range val.MapKeys() {
			// Convert key to string
			keyStr := key.String()
			if maskingFields[keyStr] {
				val.SetMapIndex(key, reflect.ValueOf("******"))
			} else {
				val.SetMapIndex(key, reflect.ValueOf(d.maskData(val.MapIndex(key).Interface())))
			}
		}
	case reflect.Struct:
		// Convert to map string so it can be masking
		// Somehow struct type is immutable
		var data map[string]interface{}
		byData, _ := json.Marshal(input)
		json.Unmarshal(byData, &data)

		// Replace the input data with new data type
		input = data
		d.maskData(data)
	}

	return input
}
