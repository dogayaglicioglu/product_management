package handler

import (
	"context"
	"log"

	"github.com/fluent/fluent-logger-golang/fluent"
)

type key int

const (
	LoggerKey  key = iota
	traceIDKey key = iota
)

func LoggerFromContext(ctx context.Context) *fluent.Fluent {
	logger, ok := ctx.Value(LoggerKey).(*fluent.Fluent)
	if !ok {
		log.Fatal("The logger is not found in the context")
		return nil
	}
	return logger
}

func TraceIdFromContext(ctx context.Context) string {
	traceId, ok := ctx.Value(traceIDKey).(string)
	if !ok {
		return ""
	}
	return traceId
}
