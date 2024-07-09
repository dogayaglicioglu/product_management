package routes

import (
	"net/http"
	"product_management/handlers"
	"product_management/middleware"

	"github.com/gorilla/mux"
)

func SetUpRoutes(r *mux.Router) {
	r.Handle("/api/users", middleware.Authorization(http.HandlerFunc(handlers.GetUsers))).Methods("GET")
}
