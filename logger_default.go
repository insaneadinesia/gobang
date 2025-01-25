package logger

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/spf13/cast"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"
)

// defaultLogger is a logger implementation using zap.
type defaultLogger struct {
	zapLogger     *zap.Logger
	enableMasking bool
	maskingFields map[string]bool
}

// NewLogger creates a new logger based on provided options.
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

// Debug logs a debug message with context and fields.
func (d *defaultLogger) Debug(ctx context.Context, message string, details ...interface{}) {
	d.log(ctx, zap.DebugLevel, message, details...)
}

// Info logs an informational message with context and fields.
func (d *defaultLogger) Info(ctx context.Context, message string, details ...interface{}) {
	d.log(ctx, zap.InfoLevel, message, details...)
}

// Warn logs a warning message with context and fields.
func (d *defaultLogger) Warn(ctx context.Context, message string, details ...interface{}) {
	d.log(ctx, zap.WarnLevel, message, details...)
}

// Error logs an error message with context and fields.
func (d *defaultLogger) Error(ctx context.Context, message string, details ...interface{}) {
	d.log(ctx, zap.ErrorLevel, message, details...)
}

// Fatal logs a fatal message with context and fields and exits the program.
func (d *defaultLogger) Fatal(ctx context.Context, message string, details ...interface{}) {
	d.log(ctx, zap.FatalLevel, message, details...)
}

// Panic logs a panic message with context and fields and panics.
func (d *defaultLogger) Panic(ctx context.Context, message string, details ...interface{}) {
	d.log(ctx, zap.PanicLevel, message, details...)
}

func (d *defaultLogger) log(ctx context.Context, level zapcore.Level, message string, details ...interface{}) {
	zapLogs := []zap.Field{}

	traceContextFields := TraceContext(ctx)

	fields := d.formatToField(details...)
	zapLogs = append(zapLogs, d.formatLogs(ctx, fields...)...)
	d.zapLogger.With(traceContextFields...).Log(level, message, zapLogs...)
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

	// Detect the message data type then convert it to json if possible
	p, ok := msg.(proto.Message)
	if ok {
		b, _ := json.Marshal(p)

		var data interface{}

		// if error happened, just print the original message
		if err := json.Unmarshal(b, &data); err != nil {
			logRecord = zap.Any(key, p)
			return
		}

		logRecord = zap.Any(key, d.maskData(data))
		return
	}

	// Detect the message data type then convert it to json if possible
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
		// Convert to map string so it can be masking. Somehow struct type is immutable
		var data map[string]interface{}
		byData, _ := json.Marshal(input)
		json.Unmarshal(byData, &data)

		// Replace the input data with new data type
		input = data
		d.maskData(data)
	}

	return input
}
