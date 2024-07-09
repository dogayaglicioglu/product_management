package handlers

import (
	"fmt"
	"net/http"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	// Example of handling an HTTP request
	fmt.Println("Handling GET request for /users")

	// You can write a response to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello from GetUsers!"))

}
