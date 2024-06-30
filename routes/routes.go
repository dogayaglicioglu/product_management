package routes

import (
	"product_management/handlers"

	"github.com/gorilla/mux"
)

func SetUpRoutes(r *mux.Router) {
	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")

}
