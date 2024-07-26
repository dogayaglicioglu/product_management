package handler

import (
	"context"

	"go.uber.org/zap"
)

type key int

const (
	LoggerKey  key = iota
	TraceIDKey key = iota
)

func LoggerFromContext(ctx context.Context) *zap.SugaredLogger {
	log, ok := ctx.Value(LoggerKey).(*zap.SugaredLogger)
	if !ok {
		return zap.NewNop().Sugar()
	}
	return log
}

func TraceIDFromContext(ctx context.Context) string {
	traceID, ok := ctx.Value(TraceIDKey).(string)
	if !ok {
		return "unknown traceID"
	}
	return traceID
}
