package routes

import (
	"net/http"
	"product_management/handlers"
	"product_management/middleware"

	"github.com/gorilla/mux"
)

func SetUpRoutes(r *mux.Router) {
	//they are same now
	r.Handle("/api/showUser", middleware.Authorization(http.HandlerFunc(handlers.GetUsers))).Methods("GET")
	r.Handle("/api/showUsers", middleware.Authorization(http.HandlerFunc(handlers.GetUsers))).Methods("GET")
}
