package logger

import (
	"context"
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

func InitLog() *zap.SugaredLogger {
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

	return LoggerInst.ZapLogger
}

func (l *LogInstance) LogWithTraceId(ctx context.Context) *zap.SugaredLogger {
	traceId, ok := ctx.Value("traceID").(string)
	if !ok {
		traceId = "unknown traceId"
	}
	return l.ZapLogger.With("traceID", traceId)

}

func (l *LogInstance) Info(ctx context.Context, msg string, args ...interface{}) {
	log := l.LogWithTraceId(ctx)
	log.Infof(msg, args...)
}

func (l *LogInstance) Error(ctx context.Context, msg string, args ...interface{}) {
	log := l.LogWithTraceId(ctx)
	log.Errorf(msg, args...)
}
