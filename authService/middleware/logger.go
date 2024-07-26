package middleware

import (
	handler "auth-service/handlers"
	"auth-service/logger"
	"context"
	"net/http"

	"github.com/google/uuid"
)

func AccessLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		logger := logger.InitLog()
		traceID := uuid.New().String()
		ctx := context.WithValue(r.Context(), handler.TraceIDKey, traceID)

		ctx = context.WithValue(r.Context(), handler.LoggerKey, logger) // to add logger to the context
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
