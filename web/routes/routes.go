package routes

import (
	"product_management/handlers"

	"github.com/gorilla/mux"
)

func SetUpRoutes(r *mux.Router) {
	//r.HandleFunc("/validate", handlers.Authorization).Methods.("GET")
	r.HandleFunc("/api/users", handlers.GetUsers).Methods("GET")

}
