package middleware

import (
	"auth-service/verify"
	"fmt"
	"net/http"
)

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			fmt.Printf("Token is empty.")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		verified := verify.VerifyToken(token)
		if verified != false {
			fmt.Println("Auth successful:")
			next.ServeHTTP(w, r)
		} else {
			fmt.Println("Unauthorized")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		}
	})

}
