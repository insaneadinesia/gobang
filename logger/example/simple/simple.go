package main

import (
	"context"
	"net/http"

	"github.com/insaneadinesia/gobang/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	ctx := context.Background()

	const serviceName = "Logger Service"
	const serviceVersion = "v1.0.0"

	// Initiate your otel (optional)
	// Using otel will automate generate trace.id
	tp := trace.NewTracerProvider()
	otel.SetTracerProvider(tp)

	ctx, span := otel.Tracer(serviceName).Start(ctx, "main")
	defer span.End()

	// Initiate logger
	logger.NewLogger(logger.Option{
		IsEnable:            true,
		EnableStackTrace:    true,
		EnableMaskingFields: true,
		MaskingFields: []string{
			"password", "token",
		},
	})

	// sample request header
	reqHeaders := http.Header{}
	reqHeaders.Add("Content-Type", "application/json")
	reqHeaders.Add("X-Client-Id", serviceName)

	// sample request body
	reqBody := map[string]interface{}{
		"request_id": "6a95fbdc-c5c4-49c2-b091-42295eff46b7",
		"name":       "mamatosai",
		"password":   "1234567",
	}

	// Prepare context logger
	ctxLogger := logger.Context{
		ServiceName:    serviceName,
		ServiceVersion: serviceVersion,
		ServicePort:    9000,
		ReqMethod:      "POST",
		ReqURI:         "/test",
		ReqBody:        reqBody,
	}

	ctx = logger.InjectCtx(ctx, ctxLogger)

	logger.Log.Info(ctx, "Request Header", reqHeaders)
	logger.Log.Info(ctx, "Request Body", reqBody)
}
