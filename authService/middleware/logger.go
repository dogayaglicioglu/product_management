package middleware

import (
	handler "auth-service/handlers"
	"auth-service/logger"
	"net/http"

	"golang.org/x/net/context"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logInstance := logger.Logger()
		if logInstance == nil {
			http.Error(w, "Logger instance is nil", http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), handler.LoggerKey, logInstance.Logger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
