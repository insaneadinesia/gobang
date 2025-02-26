# Gobang - Logger
**Gobang - Logger** is a structured logging library for Go services. It leverages [Zap](https://github.com/uber-go/zap) as its core logging mechanism and adds several additional features to enhance logging capabilities.

## Table of Contents

- [Installation](#installation)
- [Features](#features)
- [Quick Start](#quick-start)
- [Log Structure](#log-structure)
  - [Common Log](#common-log)
  - [TDR Log](#tdr-log)
- [Comparison & Explanation](#comparison--explanation)

## Installation
To install Gobang - Logger, use the following command:

```bash
go get -u insaneadinesia/gobang/logger
```

## Features
- **Formatted Logs:** Ensures that logs are consistently formatted for easier parsing and analysis.
- **Masking Data:** Allows sensitive data (e.g., passwords, tokens) to be masked in logs.
- **Additional Trace & Span ID Information:** Includes trace and span IDs for better distributed tracing support.

## Quick Start
Here's a simple example to get you started with Gobang - Logger:

```go
logger.NewLogger(logger.Option{
	IsEnable:            true,
	EnableStackTrace:    true,
	EnableMaskingFields: true,
	MaskingFields: []string{
	  "password",
      "token",
      "anything-you-want-to-mask"
	},
})

ctxLogger := logger.Context{
	ServiceName:    "service-name",
	ServiceVersion: "service-version",
	ServicePort:    9000,
	ReqMethod:      "req-method",
	ReqURI:         "req-path",
}

ctx = logger.InjectCtx(ctx, ctxLogger)

logger.Log.Info(ctx, "title", "replace-with-anything")
```
For more detailed examples, see the [example directory](https://github.com/insaneadinesia/gobang/tree/master/logger/example).

## Log Structured
### Common Log
Common logs are used for logging errors, request responses to outbound services, or additional messages to aid in issue tracing. They can be placed anywhere in your code.
```json
{
  "level": "info",
  "xtime": "2025-02-26 11:30:20.582",
  "message": "Request Header",
  "trace.id": "b03b1bba60aa5e3e8c2ee0ce141b0ad8",
  "span.id": "e4fa02ad9b3bde14",
  "app_name": "Logger Service",
  "app_version": "v1.0.0",
  "app_port": 9000,
  "app_tag": "",
  "app_method": "POST",
  "app_uri": "/test",
  "message_0": "test-message-1",
  "message_1": "test-message-2"
}
```

### TDR Log
TDR (Transaction Data Record) logs are used for logging request and response data served by your service. They are typically placed in request-response middleware.
```json
{
  "level": "info",
  "xtime": "2025-02-26 12:13:21.103",
  "message": "TDR",
  "trace.id": "b03b1bba60aa5e3e8c2ee0ce141b0ad8",
  "span.id": "e4fa02ad9b3bde14",
  "app_name": "Logger Service",
  "app_version": "v1.0.0",
  "app_port": 9000,
  "app_tag": "",
  "app_method": "POST",
  "app_uri": "/test",
  "app_response_code": 200,
  "app_exec_time": "29.283917ms",
  "app_request_header": {
    "Accept": [
      "*/*"
    ],
    "Accept-Encoding": [
      "gzip, deflate, br"
    ],
    "Connection": [
      "keep-alive"
    ],
    "Content-Length": [
      "229"
    ],
    "Content-Type": [
      "application/json"
    ],
    "Postman-Token": [
      "960b60ab-444a-4d01-a8b6-e3b74eff2524"
    ],
    "User-Agent": [
      "PostmanRuntime/7.43.0"
    ]
  },
  "app_request": {
    "amount": 100000,
    "recipient": {
      "account_no": "1234567890",
      "bank_code": "014",
      "bank_name": "bca"
    },
    "username": "mamatosai",
    "wallet_uuid": "48366b95-3320-464b-8f97-fa27a75fb0be"
  },
  "app_response": {
    "message": "Request Successfully Processed"
  }
}
```

### Comparation & Explanation
| Key | Common Log | TDR Log | Description |
|---|---|---|---|
|level | &check; | &check; | The level of the log. Possible value: debug, info, warn, error, fatal or panic. |
| x-time | &check; | &check; | Date and time when the log was created. |
| trace.id | &check; | &check; | Trace ID extracted from the traceparent context. |
| span.id | &check; | &check; | Span ID extracted from the traceparent context. |
| app_name | &check; | &check; | Name of the service. |
| app_version | &check; | &check; | Version of the service. |
| app_port | &check; | &check; | Port number on which the service is running. |
| app_method | &check; | &check; | HTTP Request Method. |
| app_uri | &check; | &check; | HTTP Request URI. |
| app_response_code | &cross; | &check; | HTTP response status code. |
| app_exec_time | &cross; | &check; | Duration of executed process. |
| app_request_header | &cross; | &check; | HTTP Request headers. |
| app_request | &cross; | &check; | HTTP Request body. |
| app_response | &cross; | &check; | HTTP Request body. |
| message | &check; | &check; | Titleor Message of the log. The log should contain at least one message. |
| message_1 | &check; | &cross; | Additional message or information. |
| message_2 | &check; | &cross; | Additional message or information. |
| message_n | &check; | &cross; | Additional message or information. |