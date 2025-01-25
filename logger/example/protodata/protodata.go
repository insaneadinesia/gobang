package main

import (
	"context"

	"github.com/insaneadinesia/gobang/logger"
	"github.com/insaneadinesia/gobang/logger/example/protodata/pb"
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
			"email",
		},
	})

	// sample request body
	reqBody := &pb.Person{
		Id:    12345,
		Name:  "mamatosai",
		Email: "rachmat.adi.p@gmail.com",
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

	logger.Log.Info(ctx, "Request Body", reqBody)
}
