package middleware

import (
	"auth-service/logger"
	"context"
	"net/http"

	"github.com/google/uuid"
)

func AccessLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.New().String()
		loggerInst := logger.LoggerInstWithTraceId(traceID)

		ctx := context.WithValue(r.Context(), logger.LoggerKey, loggerInst)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
