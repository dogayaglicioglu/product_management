package logger

import (
	"context"
	"fmt"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	serviceName     = "[GO APP]"
	logFieldService = "Authentication_Service"
	logFile         = "logger/app.log"
)

type LogInstance struct {
	ZapLogger *zap.SugaredLogger
	TraceId   string
}

var (
	LoggerInst LogInstance
)

type key int

const (
	LoggerKey  key = iota
	TraceIDKey key = iota
)

func LoggerInstWithTraceId(traceId string) *LogInstance {
	return &LogInstance{
		ZapLogger: LoggerInst.ZapLogger,
		TraceId:   traceId,
	}
}
func InitLog() {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "ts"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	w := io.MultiWriter(file, os.Stdout)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(w),
		zap.DebugLevel,
	)
	fields := zap.Fields(zap.String(logFieldService, serviceName))
	options := []zap.Option{fields, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)}
	LoggerInst.ZapLogger = zap.New(core, options...).Sugar()

}

func (l *LogInstance) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.ZapLogger == nil {
		fmt.Print("Zap Logger is nil..")
		return
	}
	l.ZapLogger.With("traceID", l.TraceId).Infof(msg, args...)
}

func (l *LogInstance) Error(ctx context.Context, msg string, args ...interface{}) {
	if l.ZapLogger == nil {
		fmt.Print("Zap Logger is nil..")
		return
	}
	l.ZapLogger.With("traceID", l.TraceId).Errorf(msg, args...)
}
