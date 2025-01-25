package logger

import (
	"context"
)

// Logger is an interface for logging messages with context and fields.
type Logger interface {
	Debug(ctx context.Context, message string, fields ...interface{})
	Info(ctx context.Context, message string, fields ...interface{})
	Warn(ctx context.Context, message string, fields ...interface{})
	Error(ctx context.Context, message string, fields ...interface{})
	Fatal(ctx context.Context, message string, fields ...interface{})
	Panic(ctx context.Context, message string, fields ...interface{})
}

// Field represents a key-value pair for logging fields.
type Field struct {
	Key string
	Val interface{}
}

// ctxKeyLogger is a type for context keys to avoid collisions.
type ctxKeyLogger struct{}

// ctxKey is the context key for the logger context.
var ctxKey = ctxKeyLogger{}

// Context holds contextual information for logging.
type Context struct {
	ServiceName    string                 `json:"app_name"`
	ServiceVersion string                 `json:"app_version"`
	ServicePort    int                    `json:"app_port"`
	Tag            string                 `json:"app_tag"`
	ReqMethod      string                 `json:"app_method"`
	ReqURI         string                 `json:"app_uri"`
	ReqHeader      string                 `json:"app_request_header"`
	ReqBody        interface{}            `json:"app_request"`
	RespBody       interface{}            `json:"app_response"`
	RespCode       int                    `json:"app_response_code,omitempty"`
	Error          string                 `json:"app_error,omitempty"`
	AdditionalData map[string]interface{} `json:"app_data,omitempty"`
	RespTime       string                 `json:"app_exec_time,omitempty"`
}

// InjectCtx injects the logger context into the parent context.
func InjectCtx(parent context.Context, ctx Context) context.Context {
	if parent == nil {
		return InjectCtx(context.Background(), ctx)
	}

	return context.WithValue(parent, ctxKey, ctx)
}

// ExtractCtx extracts the logger context from the given context.
func ExtractCtx(ctx context.Context) Context {
	if ctx == nil {
		return Context{}
	}

	val, ok := ctx.Value(ctxKey).(Context)
	if !ok {
		return Context{}
	}

	return val
}
