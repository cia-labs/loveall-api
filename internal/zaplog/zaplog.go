package zaplog

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var TempLogger *log.Logger
var Logger *zap.Logger
var LoggerErr error

func init() {
	zapconfig := zap.NewProductionEncoderConfig()
	zapconfig.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(zapconfig)
	logFile, _ := os.OpenFile(fmt.Sprintf("/var/log/loveall/loveall.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// logFile, _ := os.OpenFile(fmt.Sprintf("/Users/vamsianamalamudi/github/love-all-backend/loveall.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(zapcore.NewCore(fileEncoder, writer, defaultLogLevel))
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	defer Logger.Sync()
}
