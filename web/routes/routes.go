package routes

import (
	"net/http"
	"product_management/handlers"
	"product_management/middleware"
	"product_management/repository"

	"github.com/gorilla/mux"
)

func SetUpRoutes(r *mux.Router, webRepo repository.WebRepository) {
	handlers.SetUserRepo(webRepo)
	r.Handle("/api/showUser/{id}", middleware.Authorization(http.HandlerFunc(handlers.GetUser))).Methods("GET")
	r.Handle("/api/showUsers", middleware.Authorization(http.HandlerFunc(handlers.GetUsers))).Methods("GET")
}
