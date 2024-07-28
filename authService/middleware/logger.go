package middleware

import (
	"auth-service/logger"
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func AccessLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.New().String()
		fmt.Printf("trace id %s", traceID)
		l := logger.LoggerInstWithTraceId(traceID)

		ctx := context.WithValue(r.Context(), logger.LoggerKey, l)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
