package logger

import (
	"context"
)

type Logger interface {
	Debug(ctx context.Context, message string, fields ...interface{})
	Info(ctx context.Context, message string, fields ...interface{})
	Warn(ctx context.Context, message string, fields ...interface{})
	Error(ctx context.Context, message string, fields ...interface{})
	Fatal(ctx context.Context, message string, fields ...interface{})
	Panic(ctx context.Context, message string, fields ...interface{})
}

type Field struct {
	Key string
	Val interface{}
}

type ctxKeyLogger struct{}

var ctxKey = ctxKeyLogger{}

type Context struct {
	ServiceName    string                 `json:"_app_name"`
	ServiceVersion string                 `json:"_app_version"`
	ServicePort    int                    `json:"_app_port"`
	Tag            string                 `json:"_app_tag"`
	ReqMethod      string                 `json:"_app_method"`
	ReqURI         string                 `json:"_app_uri"`
	ReqHeader      string                 `json:"_app_request_header"`
	ReqBody        interface{}            `json:"_app_request"`
	RespBody       interface{}            `json:"_app_response"`
	RespCode       int                    `json:"_app_response_code,omitempty"`
	Error          string                 `json:"_app_error,omitempty"`
	AdditionalData map[string]interface{} `json:"_app_data,omitempty"`
	RespTime       string                 `json:"_app_exec_time,omitempty"`
}

func InjectCtx(parent context.Context, ctx Context) context.Context {
	if parent == nil {
		return InjectCtx(context.Background(), ctx)
	}

	return context.WithValue(parent, ctxKey, ctx)
}

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
