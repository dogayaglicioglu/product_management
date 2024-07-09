package middleware

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// middleware
func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			fmt.Printf("Token is empty.")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return

		}

		resp, err := http.Get("http://gateway/auth/verify?token=" + token)
		if err != nil || resp.StatusCode != http.StatusOK {
			fmt.Printf("Error in verify token func. v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		fmt.Println("Auth successful:", string(body))
		next.ServeHTTP(w, r)
	})

}
