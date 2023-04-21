package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(level string) *zap.Logger {
	var zapLevel zapcore.Level
	if err := zapLevel.Set(level); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	config := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "msg",
			LevelKey:   "level",
			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,
		},
	}

	logger, _ := config.Build()
	return logger
}

func With(fields ...zap.Field) *zap.Logger {
	return zap.L().With(fields...)
}
