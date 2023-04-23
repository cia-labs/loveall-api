package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(level string) *zap.Logger {
	var zapLevel zapcore.Level
	if err := zapLevel.Set(level); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	// config := zap.Config{
	// 	Encoding:         "json",
	// 	Level:            zap.NewAtomicLevelAt(zapLevel),
	// 	OutputPaths:      []string{"stdout"},
	// 	ErrorOutputPaths: []string{"stderr"},
	// 	EncoderConfig: zapcore.EncoderConfig{
	// 		MessageKey: "msg",
	// 		LevelKey:   "level",
	// 		TimeKey:    "time",
	// 		EncodeTime: zapcore.ISO8601TimeEncoder,
	// 	},
	// }
	logFile, _ := os.OpenFile(fmt.Sprintf("/var/log/loveall/loveall.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)
	zapconfig := zap.NewProductionEncoderConfig()
	zapconfig.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(zapconfig)
	core := zapcore.NewTee(zapcore.NewCore(fileEncoder, writer, zapLevel))
	// core.

	// logger, _ := config.Build()
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger
}

func With(fields ...zap.Field) *zap.Logger {
	return zap.L().With(fields...)
}
