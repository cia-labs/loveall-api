package zaplog

import (
	"fmt"
	"log"

	"go.uber.org/zap"
)

var TempLogger *log.Logger
var Logger *zap.Logger
var LoggerErr error

func init() {
	TempLogger = log.Default()

	Logger, LoggerErr = zap.NewProduction()
	fmt.Println(Logger)
}
