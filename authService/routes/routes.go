package routes

import (
	handler "auth-service/handlers"

	"github.com/gorilla/mux"
)

func SetUpRoutes(r *mux.Router) {
	r.HandleFunc("/auth/register", handler.Register).Methods("POST")
	r.HandleFunc("/auth/login", handler.Login).Methods("POST")
	r.HandleFunc("/auth/verify", handler.Verify).Methods("GET")
}
