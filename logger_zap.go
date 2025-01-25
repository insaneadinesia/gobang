package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapLogger creates a new zap logger based on the provided options.
func NewZapLogger(opt Option) *zap.Logger {
	var outputPaths, errorOutputPaths []string

	// Default path will be shown on console
	if opt.IsEnable {
		outputPaths = []string{"stdout"}
		errorOutputPaths = []string{"stderr"}
	}

	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       true,
		DisableCaller:     true,
		DisableStacktrace: !opt.EnableStackTrace,
		Encoding:          "json",
		EncoderConfig:     getEncoderConfig(),
		OutputPaths:       outputPaths,
		ErrorOutputPaths:  errorOutputPaths,
	}

	return zap.Must(cfg.Build())
}

// getEncoderConfig returns the encoder configuration for zap.
func getEncoderConfig() zapcore.EncoderConfig {
	cfg := zap.NewProductionEncoderConfig()
	cfg.TimeKey = "xtime"
	cfg.EncodeDuration = zapcore.MillisDurationEncoder
	cfg.EncodeTime = timeEncoder
	cfg.MessageKey = "message"

	return cfg
}

// timeEncoder formats the time for zap logs.
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.999"))
}
