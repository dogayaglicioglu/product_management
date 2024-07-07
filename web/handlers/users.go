package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func Authorization(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return

	}

	resp, err := http.Get("http://gateway/auth/verify/?token=" + token)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	fmt.Println("Auth successful:", string(body))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Auth successful!"))

}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	// Example of handling an HTTP request
	fmt.Println("Handling GET request for /users")

	// You can write a response to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello from GetUsers!"))

}
